/*
Copyright © 2022 Anttu Suhonen

*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var city string
var apiKey string
var units string
var forecast string
var output string
var debug bool = false

var rootCmd = &cobra.Command{
	Use:   "wthroo",
	Short: "A command line OpenWeatherMap lookup tool.",
	Long:  `A simple command line OpenWeatherMap lookup tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		if city == "" || apiKey == "" {
			cmd.Usage()
			os.Exit(1)
		}

		if debug {
			fmt.Println("City: ", city)
			os.Exit(0)
		}

		weather := getWeather(city, units, apiKey, forecast, output)
		buildMessage(weather)

	},
}

func buildMessage(weather []byte) {
	fmt.Printf("%s", weather)
}

func getWeather(city, units, apiKey, forecast, output string) []byte {
	request := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&units=%s&appid=%s", city, units, apiKey)

	response, err := http.Get(request)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data bytes.Buffer
	error := json.Indent(&data, responseData, "", "  ")
	if error != nil {
		log.Println("JSON parse error: ", error)
		os.Exit(1)
	}

	return data.Bytes()
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// pointer, name, shorthand, default value, info
	viper.AutomaticEnv()
	flags := rootCmd.Flags()
	flags.StringVarP(&apiKey, "apiKey", "k", "", "API key from OpenWeatherMap.")
	viper.BindPFlag("apiKey", flags.Lookup("apiKey"))
	flags.StringVarP(&city, "city", "c", "", "City to get weather info from. Can be city name, state code and country code divided by comma.")
	viper.BindPFlag("city", flags.Lookup("city"))
	flags.StringVarP(&units, "units", "u", "metric", "Optional, default metric. Units of measurement. Standard, metric and imperial units are available.")
	viper.BindPFlag("apiKey", flags.Lookup("apiKey"))
	flags.StringVarP(&forecast, "forecast", "f", "", "Optional. Type of forecast, default is current values. Options: hourly, daily, climate.")
	viper.BindPFlag("forecast", flags.Lookup("forecast"))
	flags.StringVarP(&output, "output", "o", "json", "Optional. Data format. Possible values are json and xml. If the mode parameter is empty the format is JSON by default.")
	viper.BindPFlag("output", flags.Lookup("output"))

	if os.Getenv("DEBUG") != "" {
		debug = true
	}
}
