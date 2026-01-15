package playerLogic

type NPC struct{
	DialogueLines map[string][]string `json:"dialogue_lines"`
	Player
}

func NewNPC(name string, description string) *NPC{
	return &NPC{
		Player: Player{
			Name: name,
			Description: description,
		},
		DialogueLines: make(map[string][]string),
	}
}

func NewNPCWithAttributes(name string, description string, attributes PlayerAttributes) *NPC{
	npc := NewNPC(name, description)
	npc.Attributes = attributes
	npc.SetAttributeModifiers()
	npc.SetDerivedStats()
	return npc
}

func (npc *NPC) AddDialogue(topic string, lines []string){
	npc.DialogueLines[topic] = lines
}

func (npc *NPC) EditDialogue(topic string, line string, index int){
	if existingLines, exists := npc.DialogueLines[topic]; exists{
		if index >= 0 && index < len(existingLines){
			existingLines[index] = line
			npc.DialogueLines[topic] = existingLines
		}
	}
}

func (npc *NPC) RemoveDialogue(topic string){
	delete(npc.DialogueLines, topic)
}

func (npc *NPC) GetDialogue(topic string) []string{
	if lines, exists := npc.DialogueLines[topic]; exists{
		return lines
	}
	return []string{}
}