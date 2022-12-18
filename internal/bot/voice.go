package bot

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

type VoiceContext struct {
	stop    chan error
	s       *discordgo.Session
	playing sync.RWMutex
}

func (vc *VoiceContext) Play(reader io.ReadCloser, guildId string, channelId string) error {
	if vc.playing.TryLock() {
		go func() {
			defer vc.playing.Unlock()
			defer reader.Close()

			voiceConnection, err := vc.s.ChannelVoiceJoin(guildId, channelId, false, true)
			if err != nil {
				fmt.Println("Bot Voice Play: Error - ", err)
				return
			}
			defer voiceConnection.Disconnect()

			time.Sleep(250 * time.Millisecond)
			defer time.Sleep(250 * time.Millisecond)
			voiceConnection.Speaking(true)
			defer voiceConnection.Speaking(false)

			opts := dca.StdEncodeOptions

			encodeSession, err := dca.EncodeMem(reader, opts)
			if err != nil {
				fmt.Println("Bot Voice Play: Error - ", err)
				return
			}
			defer encodeSession.Cleanup()

			done := make(chan error)
			stream := dca.NewStream(encodeSession, voiceConnection, done)
			ticker := time.NewTicker(5 * time.Second)

			for {
				select {
				case err := <-vc.stop:
					if err != nil {
						fmt.Println("Bot Voice Play: Error - ", err)
					}
					fmt.Println("Bot Voice Play: Stream stopped")
					return
				case err := <-done:
					if err != nil && err != io.EOF {
						fmt.Println("Bot Voice Play: Error - ", err)
					}
					fmt.Println("Bot Voice Play: Stream done")
					return
				case <-ticker.C:
					stats := encodeSession.Stats()
					playbackPosition := stream.PlaybackPosition()

					fmt.Printf("Bot Voice Play: Playback: %10s, Transcode Stats: Time: %5s, Size: %5dkB, Bitrate: %6.2fkB, Speed: %5.1fx\n", playbackPosition, stats.Duration.String(), stats.Size, stats.Bitrate, stats.Speed)
				}
			}
		}()
	} else {
		reader.Close()
		return errors.New("Already playing, stop first")
	}
	return nil
}

func (vc *VoiceContext) Stop() error {
	if vc.playing.TryLock() {
		vc.playing.Unlock()
		return errors.New("Not playing")
	}
	vc.stop <- nil
	return nil
}
