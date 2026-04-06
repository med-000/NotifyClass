package discord

import (
	"fmt"
	"strconv"
	"strings"
)

const slotCommandLength = 11

func IsSlotCommand(content string) bool {
	trimmed := strings.TrimSpace(content)
	return len(trimmed) == slotCommandLength && strings.HasPrefix(trimmed, "/")
}

func ParseSlotCommand(content string) (SlotQuery, error) {
	trimmed := strings.TrimSpace(content)
	if !IsSlotCommand(trimmed) {
		return SlotQuery{}, fmt.Errorf("invalid slot command format")
	}

	body := trimmed[1:]

	year, err := strconv.Atoi(body[0:4])
	if err != nil {
		return SlotQuery{}, fmt.Errorf("invalid year: %w", err)
	}

	term, err := strconv.Atoi(body[4:6])
	if err != nil {
		return SlotQuery{}, fmt.Errorf("invalid term: %w", err)
	}

	day, err := strconv.Atoi(body[6:8])
	if err != nil {
		return SlotQuery{}, fmt.Errorf("invalid day: %w", err)
	}

	period, err := strconv.Atoi(body[8:10])
	if err != nil {
		return SlotQuery{}, fmt.Errorf("invalid period: %w", err)
	}

	if term < 1 || term > 2 {
		return SlotQuery{}, fmt.Errorf("term out of range: %d", term)
	}
	if day < 1 || day > 7 {
		return SlotQuery{}, fmt.Errorf("day out of range: %d", day)
	}
	if period < 1 {
		return SlotQuery{}, fmt.Errorf("period out of range: %d", period)
	}

	return SlotQuery{
		Year:   year,
		Term:   term,
		Day:    day,
		Period: period,
	}, nil
}
