package exceptionist

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func parseFile(filepath string) bucket {

	viper.SetConfigName("messages_tr.properties")
	viper.AddConfigPath(os.Getenv("GOPATH") + "/src/mytest/messages")
	//viper.AddConfigPath(os.Getenv("CONFIG_PATH"))
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error messages file: %s", err))
	}

	var bucket bucket = make(bucket)
	err = viper.Unmarshal(bucket)
	if err != nil {
		panic(fmt.Errorf("Fatal error messages file: %s", err))
	}
	//_, err := toml.DecodeFile(filepath, &bucket)
	//if err != nil {
	//	log.Fatal(err)
	//}

	return bucket
}
