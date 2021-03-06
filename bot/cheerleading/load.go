package cheerleading

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

const (
	rootFileDir  = "./resource/"
	rootFileName = "voicebanks.json"
)

func init() {

	rootDecoded := []string{}

	log.Println("load voicebank files...")

	// parse json file
	if f, err := os.Open(rootFileDir + rootFileName); err != nil {
		log.Fatalln("cannot open route file: " + err.Error())
	} else {

		defer f.Close()
		if err := json.NewDecoder(f).Decode(&rootDecoded); err != nil {
			log.Fatalln("cannot parse route json file: " + err.Error())
		}
	}

	for _, relPath := range rootDecoded {

		// parse voicebank file (json)
		func() {
			if f, err := os.Open(rootFileDir + relPath); err != nil {
				log.Fatalln("cannot open voicebank file (path: " + relPath + "): " + err.Error())
			} else {

				defer f.Close()
				voicebank := Cheerleader{}
				if err := json.NewDecoder(f).Decode(&voicebank); err != nil {
					log.Fatalln("cannot open voicebank file (path: " + relPath + "): " + err.Error())
				} else {
					for i := range voicebank.Talkset {
						for j := range voicebank.Talkset[i] {
							voicebank.Talkset[i][j].FileName = rootFileDir + filepath.Dir(relPath) + "/" + voicebank.Talkset[i][j].FileName
						}
					}
					Cheerleaders = append(Cheerleaders, voicebank)
				}
			}
		}()
	}

	log.Println("load voicebank files ended")
}
