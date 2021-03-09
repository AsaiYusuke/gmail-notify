package main

import "text/template"

type commandTemplate struct {
	createTemplate *template.Template
	removeTemplate *template.Template
	countTemplate  *template.Template
}

func (t *commandTemplate) init() {
	t.createTemplate = template.New(commandTemplateCreateTemplateID)
	t.createTemplate.Parse(commandTemplateCreateTemplate)
	t.removeTemplate = template.New(commandTemplateRemoveTemplateID)
	t.removeTemplate.Parse(commandTemplateRemoveTemplate)
	t.countTemplate = template.New(commandTemplateCountTemplateID)
	t.countTemplate.Parse(commandTemplateCountTemplate)
}
