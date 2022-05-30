package main

type JambonzSay struct {
	Verb        string
	Text        string
	Loop        bool
	EarlyMedia  bool
	Synthesizer struct {
		Vendor   string
		Language string
		Gender   string
		Voice    string
	}
}
