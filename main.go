//BY GRANDPA HACKER 2023

package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	gpt3 "github.com/PullRequestInc/go-gpt3"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

type NullWriter int

func (NullWriter) Write([]byte) (int, error) { return 0, nil }

func main() {
	log.SetOutput(new(NullWriter))
	apiKey := goDotEnvVariable("API_KEY")

	if apiKey == "" {
		panic("Error loading api key")
	}

	ctx := context.Background()
	client := gpt3.NewClient(apiKey)
	rootCmd := &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			scanner := bufio.NewScanner(os.Stdin)
			quit := false

			//quit loop
			for !quit {
				fmt.Print("Вкажіть ваше питання (q шоб вийти):")

				if !scanner.Scan() {
					break
				}

				question := scanner.Text()
				questionParam := validateQuestion(question)
				switch questionParam {
				case "q":
					quit = true
				case "":
					continue

				default:
					GetResponse(client, ctx, questionParam)
				}
			}
		},
	}

	log.Fatal(rootCmd.Execute())
}

// .ENV file reader
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func GetResponse(client gpt3.Client, ctx context.Context, quesiton string) {
	err := client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			quesiton,
		},
		MaxTokens:   gpt3.IntPtr(3000),
		Temperature: gpt3.Float32Ptr(0),
	}, func(resp *gpt3.CompletionResponse) {
		fmt.Print(resp.Choices[0].Text)
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(13)
	}
	fmt.Printf("\n")
}

func validateQuestion(question string) string {
	quest := strings.Trim(question, " ")
	keywords := []string{"", "loop", "break", "continue", "cls", "exit", "block"}
	for _, x := range keywords {
		if quest == x {
			return ""
		}
	}
	return quest
}
