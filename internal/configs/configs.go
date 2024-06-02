package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	Port                  int    `mapstructure:"PORT"`
	JWTSecretKey          string `mapstructure:"JWT_SECRET_KEY"`
	MongoUri              string `mapstructure:"MONGODB_URI"`
	Environment           string `mapstructure:"ENVIRONMENT"`
	WhatsappToken         string `mapstructure:"WHATSAPP_TOKEN"`
	FacebookGraphEndpoint string `mapstucture:"FACEBOOK_GRAPH_ENDPOINT"`
	PhoneNumberID         string `mapstucture:"PHONE_NUMBER_ID"`
	RedisUri              string `mapstucture:"REDIS_URI"`
	ATSMSEndpoint         string `mapstucture:"AT_SMS_ENDPOINT"`
	ATAPIKey              string `mapstucture:"AT_API_KEY"`
	ShortCode             string `mapstucture:"SHORTCODE"`
	Username              string `mapstucture:"USERNAME"`
}

func LoadConfig() (config Config, err error) {
	port, jwtSecretKey, mongoUri, environment, whatsappToken, facebookGraphEndpoint, phonenUmberID, redisUri, ATAPIKey, ATSMSEndpoint, shortCode, username :=
		os.Getenv("PORT"), os.Getenv("JWT_SECRET_KEY"), os.Getenv("MONGODB_URI"),
		os.Getenv("ENVIRONMENT"), os.Getenv("WHATSAPP_TOKEN"), os.Getenv("FACEBOOK_GRAPH_ENDPOINT"),
		os.Getenv("PHONE_NUMBER_ID"), os.Getenv("REDIS_URI"), os.Getenv("AT_API_KEY"), os.Getenv("AT_SMS_ENDPOINT"),
		os.Getenv("SHORTCODE"), os.Getenv("USERNAME")

	if whatsappToken != "" && redisUri != "" && mongoUri != "" && facebookGraphEndpoint != "" && phonenUmberID != "" {
		config.Port, _ = strconv.Atoi(port)
		config.JWTSecretKey = jwtSecretKey
		config.MongoUri = mongoUri
		config.Environment = environment
		config.WhatsappToken = whatsappToken
		config.FacebookGraphEndpoint = facebookGraphEndpoint
		config.PhoneNumberID = phonenUmberID
		config.RedisUri = redisUri
		config.ATAPIKey = ATAPIKey
		config.ATSMSEndpoint = ATSMSEndpoint
		config.ShortCode = shortCode
		config.Username = username

		return config, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current working directory:", err)
	}

	configPath := filepath.Join(cwd, "./")

	viper.AddConfigPath(configPath)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	return config, nil
}
