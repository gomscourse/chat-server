package chat

import (
	"context"
	"fmt"
)

func (s chatService) SendMessage(ctx context.Context, sender, text string, chatID int64) error {
	id, err := s.repo.CreateMessage(ctx, chatID, sender, text)
	if err != nil {
		return err
	}
	fmt.Println("created message with id", id)

	//send notification and similar staff

	return nil
}
