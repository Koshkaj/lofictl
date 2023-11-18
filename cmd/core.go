package cmd

import (
	"errors"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/wav"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func processFile(fileName string, ratio float64, boost bool) error {
	fp := strings.Split(fileName, ".")
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return err
	}
	defer streamer.Close()
	var mainStreamer beep.Streamer

	if boost {
		mainStreamer = &effects.Gain{Gain: 1.3, Streamer: streamer}
	} else {
		mainStreamer = streamer
	}

	reSample := beep.ResampleRatio(4, ratio, mainStreamer)
	fileOut, err := os.Create(fmt.Sprintf("%s_out.wav", fp[0]))
	if err != nil {
		return err
	}
	defer fileOut.Close()
	err = wav.Encode(fileOut, reSample, beep.Format{
		SampleRate:  format.SampleRate,
		NumChannels: 2,
		Precision:   2,
	})
	if err != nil {
		return err
	}
	return nil
}

func applyCommand(cmd *cobra.Command, args []string) error {
	fileName, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}
	ratio, err := cmd.Flags().GetFloat64("ratio")
	if err != nil {
		return err
	}
	boost, err := cmd.Flags().GetBool("boost")
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		var wg sync.WaitGroup
		// handle directory parameter
		// call process files in goroutines
		filepath.Walk(fileName, func(fileInfo string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(fileInfo, ".mp3") {
				// spawn 5 goroutines that would process these files and ouput to a directory afterwards
				wg.Add(1)
				go func() {
					defer wg.Done()
					_ = processFile(fileInfo, ratio, boost)
				}()
			}
			return nil

		})
		wg.Wait()
	} else {
		// else just return a single file handler
		if !strings.HasSuffix(fileName, ".mp3") {
			return errors.New("file is not an mp3 extension")
		}
		err = processFile(fileName, ratio, boost)
		return err
	}

	return nil
}
