package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	// 提供操作选项
	fmt.Println("请选择操作类型:")
	fmt.Println("1: 备份应用数据")
	fmt.Println("2: 恢复应用数据")

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("请输入选项序号: ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var packageName, localBackupPath, backupFilePath string

	if choice == "1" {
		fmt.Print("请输入应用包名: ")
		packageName, _ = reader.ReadString('\n')
		packageName = strings.TrimSpace(packageName)

		fmt.Print("请输入本地备份路径 (默认: ./backup_data): ")
		localBackupPath, _ = reader.ReadString('\n')
		localBackupPath = strings.TrimSpace(localBackupPath)
		if localBackupPath == "" {
			localBackupPath = "./backup_data"
		}

		backupAppData(packageName, localBackupPath)
	} else if choice == "2" {
		fmt.Print("请输入应用包名: ")
		packageName, _ = reader.ReadString('\n')
		packageName = strings.TrimSpace(packageName)

		fmt.Print("请输入备份文件(绝对路径)路径: ")
		backupFilePath, _ = reader.ReadString('\n')
		backupFilePath = strings.TrimSpace(backupFilePath)

		if backupFilePath == "" {
			fmt.Println("备份文件路径不能为空")
			return
		}

		restoreAppData(packageName, backupFilePath)
	} else {
		fmt.Println("无效的选项，请选择 1 或 2")
	}
}

func backupAppData(packageName, localBackupPath string) {
	// 根据包名和具体时间生成备份文件名
	timestamp := time.Now().Format("20060102_150405")
	deviceBackupPath := fmt.Sprintf("/sdcard/%s_backup_%s.tar.gz", packageName, timestamp)

	// 确保本地备份路径存在
	if err := os.MkdirAll(localBackupPath, os.ModePerm); err != nil {
		fmt.Printf("创建本地备份路径失败: %v\n", err)
		return
	}

	// 定义备份命令
	backupCommands := []string{
		fmt.Sprintf("am force-stop %s", packageName),
		fmt.Sprintf("rm -rf /data/data/.external.%s", packageName),
		fmt.Sprintf("ln -sf \"/sdcard/Android/data/%s\" \"/data/data/.external.%s\"", packageName, packageName),
		fmt.Sprintf(
			"tar -c \"/data/data/%s/.\" \"/data/data/.external.%s/.\" --exclude \"data/%s/./lib*\" --exclude \"data/data/%s/./cache\" --exclude \"data/data/%s/./dex\" --exclude \"data/data/%s/./app_ras_blobs\" --exclude \"data/data/%s/./app_msqrd*\" | gzip > \"%s\"",
			packageName, packageName, packageName, packageName, packageName, packageName, packageName, deviceBackupPath,
		),
		fmt.Sprintf("rm -rf \"/data/data/.external.%s\"", packageName),
		fmt.Sprintf("am force-stop %s", packageName),
		fmt.Sprintf("chown media_rw:media_rw \"%s\"", deviceBackupPath),
	}

	// 使用 && 将命令拼接为单行
	fullCommand := strings.Join(backupCommands, " && ")

	// 执行完整命令
	if err := executeADBCommand(fmt.Sprintf("su -c \"%s\"", fullCommand)); err != nil {
		fmt.Printf("备份失败: %v\n", err)
		return
	}

	// 将备份文件从设备复制到本地
	if err := pullBackupFile(deviceBackupPath, localBackupPath); err != nil {
		fmt.Printf("拉取备份文件失败: %v\n", err)
		return
	}

	fmt.Println(localBackupPath, "备份完成!")
}

func restoreAppData(packageName, backupFilePath string) {
	// 将备份文件推送到设备
	deviceBackupPath := fmt.Sprintf("/sdcard/%s_restore.tar.gz", packageName)
	if err := pushBackupFile(backupFilePath, deviceBackupPath); err != nil {
		fmt.Printf("推送备份文件失败: %v\n", err)
		return
	}

	// 获取应用的 userId
	userIdCommand := fmt.Sprintf("dumpsys package %s | grep userId=", packageName)
	output, err := executeADBCommandWithOutput(userIdCommand)
	if err != nil {
		fmt.Printf("获取 userId 失败: %v\n", err)
		return
	}

	// 解析 userId
	userId := strings.TrimSpace(strings.Split(output, "=")[1])

	// 定义恢复命令
	restoreCommands := []string{
		fmt.Sprintf("pm clear %s", packageName),
		fmt.Sprintf("mkdir -p /sdcard/Android/data/%s/files", packageName),
		fmt.Sprintf("rm -Rf \"/data/data/.external.%s\"", packageName),
		fmt.Sprintf("ln -sf \"/sdcard/Android/data/%s\" \"/data/data/.external.%s\"", packageName, packageName),
		fmt.Sprintf("cat \"%s\" | gunzip | tar -C \"/\" -x --exclude data/data/%s/lib --exclude data/data/%s/./lib* --exclude \"data/data/%s/./dex\" --exclude \"data/data/%s/./app_ras_blobs\" --exclude \"data/data/%s/./app_msqrd*\"", deviceBackupPath, packageName, packageName, packageName, packageName, packageName),
		fmt.Sprintf("rm -Rf \"/data/data/.external.%s\"", packageName),
		fmt.Sprintf("chown -R media_rw:media_rw \"/sdcard/Android/data/%s\"", packageName),
		fmt.Sprintf("chown -hR %s:%s \"/data/data/%s\"", userId, userId, packageName),
		fmt.Sprintf("chmod -R u+rwx \"/data/data/%s\"", packageName),
		fmt.Sprintf("/system/bin/restorecon -R \"/data/data/%s\"", packageName),
		fmt.Sprintf("rm -Rf \"/data/data/.%s\"", packageName),
	}

	// 使用 && 将命令拼接为单行
	fullCommand := strings.Join(restoreCommands, " && ")

	// 执行完整命令
	if err := executeADBCommand(fmt.Sprintf("su -c \"%s\"", fullCommand)); err != nil {
		fmt.Printf("恢复失败: %v\n", err)
		return
	}

	fmt.Println("恢复完成!")
}

func executeADBCommand(command string) error {
	cmd := exec.Command("adb", "shell", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func executeADBCommandWithOutput(command string) (string, error) {
	cmd := exec.Command("adb", "shell", command)
	output, err := cmd.Output()
	return string(output), err
}

func pullBackupFile(devicePath, localPath string) error {
	cmd := exec.Command("adb", "pull", devicePath, localPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func pushBackupFile(localPath, devicePath string) error {
	cmd := exec.Command("adb", "push", localPath, devicePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
