package sworld

import "time"

// SkillSource represents a source for a skill
type SkillSource interface {
	Damage() int
}

// SkillTarget represents the target for a skill
type SkillTarget interface {
	ReceiveDamage(source Skill, amount int) int
}

// Skill represents a skill
type Skill interface {
	Use(SkillTarget) error
	WaitTime() time.Duration
}

// HitSkill is a basic skill
type HitSkill struct {
	lastUse  time.Time
	source   SkillSource
	cooldown time.Duration
}

// NewHitSkill creates a new HitSkill
func NewHitSkill(source SkillSource) *HitSkill {
	return &HitSkill{
		cooldown: time.Millisecond * 500,
		source:   source,
	}
}

// WaitTime is the time before this skill can be used
func (h *HitSkill) WaitTime() time.Duration {
	if h.lastUse.IsZero() {
		return 0
	}

	elapsed := time.Since(h.lastUse)
	if elapsed >= h.cooldown {
		return 0
	}
	return h.cooldown - elapsed
}

// Use the skill againgst a target
func (h *HitSkill) Use(target SkillTarget) error {
	damage := h.source.Damage()

	target.ReceiveDamage(h, damage)

	h.lastUse = time.Now()
	return nil
}
