package main

import (
	"encoding/base64"
	"fmt"

	webpush "github.com/SherClockHolmes/webpush-go"
)

func main() {
	privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		panic(err)
	}

	fmt.Println("VAPID Keys Generated Successfully!")
	fmt.Println("=====================================")
	fmt.Println()
	fmt.Println("Add these to your .env file:")
	fmt.Println()
	fmt.Printf("VAPID_PUBLIC_KEY=%s\n", publicKey)
	fmt.Printf("VAPID_PRIVATE_KEY=%s\n", privateKey)
	fmt.Println("VAPID_SUBJECT=mailto:your-email@example.com")
	fmt.Println()
	fmt.Println("Make sure to replace 'your-email@example.com' with your actual email address.")
	fmt.Println()
	fmt.Println("Public Key (for frontend):", publicKey)
	fmt.Println()
	fmt.Println("Keys are base64 encoded and ready to use.")
}

func init() {
	// Verify the keys can be decoded
	testPrivate, testPublic, _ := webpush.GenerateVAPIDKeys()
	_, err1 := base64.RawURLEncoding.DecodeString(testPrivate)
	_, err2 := base64.RawURLEncoding.DecodeString(testPublic)
	if err1 != nil || err2 != nil {
		panic("Failed to generate valid VAPID keys")
	}
}
