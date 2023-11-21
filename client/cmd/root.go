/*
Copyright © 2023 NAME HERE abhinavmohan12@gmail.com
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// url = "https://example.com/api/endpoint"
type Command struct {
	Command  string `json:"command"`
	Filename string `json:"filename"`
}

var pwd = "/"

var rootCmd = &cobra.Command{
	Use:   "panda",
	Short: "This is the root command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Executing root command")
	},
}
var touchCmd = &cobra.Command{
	Use:   "touch [filename]",
	Short: "remove file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		url := "http://localhost:9092/tasks"

		// json_string := "{\"cmd\": touch, \"filename\":" + filename + ", pwd:" + pwd + "}"
		json_string := map[string]interface{}{
			"cmd":      "touch",
			"filename": filename,
			"pwd":      pwd,
		}
		fmt.Print(json_string)
		json_string1, err := json.Marshal(json_string)

		if err != nil {
			fmt.Println("Erron encoding", err)
			return
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_string1))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			fmt.Println("Request successful", http.StatusOK)
		} else {
			fmt.Println("Request failed with status code:", resp.StatusCode)
		}
	},
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list all files",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		url := "http://localhost:9092/tasks"

		// json_string := "{Command: ls}"
		json_string := map[string]interface{}{
			"cmd": "ls",
		}
		json_string1, err := json.Marshal(json_string)

		if err != nil {
			fmt.Println("Erron encoding", err)
			return
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_string1))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			fmt.Println("Request successful", http.StatusOK)
		} else {
			fmt.Println("Request failed with status code:", resp.StatusCode)
		}
	},
}

var rmCmd = &cobra.Command{
	Use:   "rm [filename]",
	Short: "remove file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		url := "http://localhost:9092/tasks"

		// json_string := "{Command: rm, Filename:" + filename + ", pwd:" + pwd + "}"
		json_string := map[string]interface{}{
			"cmd":      "rm",
			"filename": filename,
			"pwd":      pwd,
		}
		json_string1, err := json.Marshal(json_string)

		if err != nil {
			fmt.Println("Erron encoding", err)
			return
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_string1))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			fmt.Println("Request successful", http.StatusOK)
		} else {
			fmt.Println("Request failed with status code:", resp.StatusCode)
		}
	},
}

type Move struct {
	Command   string `json:"command"`
	Filename1 string `json:"filename1"`
	Filename2 string `json:"filename2"`
}

var mvCmd = &cobra.Command{
	Use:   "mv",
	Short: "move or rename files",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		filename1 := args[0]
		filename2 := args[1]
		fmt.Print(filename1, filename2)
		url := "http://localhost:9092/tasks"

		// json_string := "{Command: mv, Filename1:" + filename1 + ", Filename2: " + filename2 + ", pwd:" + pwd + "}"
		json_string := map[string]interface{}{
			"cmd":       "mv",
			"filename1": filename1,
			"filename2": filename2,
			"pwd":       pwd,
		}
		json_string1, err := json.Marshal(json_string)
		fmt.Print(json_string)
		if err != nil {
			fmt.Println("Erron encoding", err)
			return
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_string1))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			fmt.Println("Request successful", http.StatusOK)
		} else {
			fmt.Println("Request failed with status code:", resp.StatusCode)
		}
	},
}

var cpCmd = &cobra.Command{
	Use:   "cp",
	Short: "copy files",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		filename1 := args[0]
		filename2 := args[1]
		url := "http://localhost:9092/tasks"

		// json_string := "{Command: cp, Filename1:" + filename1 + ", Filename2: " + filename2 + ", pwd:" + pwd + "}"
		json_string := map[string]interface{}{
			"cmd":       "cp",
			"filename1": filename1,
			"filename2": filename2,
			"pwd":       pwd,
		}
		json_string1, err := json.Marshal(json_string)
		fmt.Print(json_string)
		if err != nil {
			fmt.Println("Erron encoding", err)
			return
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_string1))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			fmt.Println("Request successful", http.StatusOK)
		} else {
			fmt.Println("Request failed with status code:", resp.StatusCode)
		}
	},
}

var cdCmd = &cobra.Command{
	Use:   "cd [filename]",
	Short: "change directory",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		url := "http://localhost:9092/tasks"
		if filename != ".." {
			pwd = filename
		}
		// json_string := "{Command: cd, Filename:" + filename + ", pwd:" + pwd + "}"
		json_string := map[string]interface{}{
			"cmd":      "cd",
			"filename": filename,
			"pwd":      pwd,
		}
		json_string1, err := json.Marshal(json_string)

		if err != nil {
			fmt.Println("Erron encoding", err)
			return
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_string1))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			fmt.Println("Request successful", http.StatusOK)
		} else {
			fmt.Println("Request failed with status code:", resp.StatusCode)
		}

	},
}

var uploadCmd = &cobra.Command{
	Use:   "upload [filename]",
	Short: "upload file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		url := "http://localhost:9092/tasks"

		// json_string := "{Command: upload, Filename:" + filename + ", pwd:" + pwd + "}"
		json_string := map[string]interface{}{
			"cmd":      "upload",
			"filename": filename,
			"pwd":      pwd,
		}
		metadata, err := os.Stat("./test.txt")
		if err != nil {
			fmt.Print("Error", err)
			return
		}
		// metadata_string := "{Filename: " + metadata.Name() + " Size: " + strconv.Itoa(int(metadata.Size())) + "}"
		metadata_string := map[string]interface{}{
			"filename": metadata.Name(),
			"size":     strconv.Itoa(int(metadata.Size())),
		}
		json_string1, err := json.Marshal(json_string)
		if err != nil {
			fmt.Println("Erron encoding", err)
			return
		}
		metadata_string1, err := json.Marshal(metadata_string)
		if err != nil {
			fmt.Println("Erron encoding", err)
			return
		}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_string1))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req1, err := http.NewRequest("POST", url, bytes.NewBuffer(metadata_string1))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		req1.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		client1 := &http.Client{}
		resp1, err := client1.Do(req1)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		defer resp.Body.Close()
		defer resp1.Body.Close()

		if resp.StatusCode == http.StatusOK {
			fmt.Println("Request successful", http.StatusOK)
		} else {
			fmt.Println("Request failed with status code:", resp.StatusCode)
		}

		url1 := "http://localhost:8000"

		// Open the file
		file, err := os.Open("./test.txt")
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		// Create a buffer to store the file content
		var buffer bytes.Buffer
		writer := multipart.NewWriter(&buffer)

		// Create a form file field and write the file content
		fileField, err := writer.CreateFormFile("file", "example.txt")
		if err != nil {
			fmt.Println("Error creating form file:", err)
			return
		}
		_, err = io.Copy(fileField, file)
		if err != nil {
			fmt.Println("Error copying file content:", err)
			return
		}

		// Close the multipart writer to finalize the request
		writer.Close()

		// Create a POST request with the multipart form data
		req2, err := http.NewRequest("POST", url1, &buffer)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// Send the request
		client2 := &http.Client{}
		resp2, err := client2.Do(req2)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		defer resp2.Body.Close()

		// Process the response as needed
		fmt.Println("Response status code:", resp.StatusCode)
	},
}

var downloadCmd = &cobra.Command{
	Use:   "download [filename]",
	Short: "dowload file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		url := "http://localhost:9092/tasks"

		// json_string := "{Command: download, Filename:" + filename + ", pwd:" + pwd + "}"
		json_string := map[string]interface{}{
			"cmd":      "download",
			"filename": filename,
			"pwd":      pwd,
		}
		json_string1, err := json.Marshal(json_string)

		if err != nil {
			fmt.Println("Erron encoding", err)
			return
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_string1))
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			fmt.Println("Request successful", http.StatusOK)
		} else {
			fmt.Println("Request failed with status code:", resp.StatusCode)
		}
	},
}

var pwdCmd = &cobra.Command{
	Use:   "pwd",
	Short: "present working directory",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		// url := "http://localhost:9092/tasks"

		json_string := map[string]interface{}{
			"cmd": "pwd",
		}
		fmt.Print(json_string)
		// json_string1, err := json.Marshal(json_string)

		// if err != nil {
		// 	fmt.Println("Erron encoding", err)
		// 	return
		// }

		// req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_string1))
		// if err != nil {
		// 	fmt.Println("Error creating request:", err)
		// 	return
		// }
		// req.Header.Set("Content-Type", "application/json")

		// client := &http.Client{}
		// resp, err := client.Do(req)
		// if err != nil {
		// 	fmt.Println("Error making request:", err)
		// 	return
		// }
		// defer resp.Body.Close()

		// if resp.StatusCode == http.StatusOK {
		// 	fmt.Println("Request successful", http.StatusOK)
		// } else {
		// 	fmt.Println("Request failed with status code:", resp.StatusCode)
		//}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bd_cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(lsCmd)
	rootCmd.AddCommand(touchCmd)
	rootCmd.AddCommand(rmCmd)
	rootCmd.AddCommand(mvCmd)
	rootCmd.AddCommand(cpCmd)
	rootCmd.AddCommand(cdCmd)
	rootCmd.AddCommand(uploadCmd)
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(pwdCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".bd_cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".bd_cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
