package manager

import debugio "virga-player/debug/io"

func (m *Manager) copy() {
	dump := m.buf.Dump()
	if err := debugio.CopyToClipboard(dump); err != nil {
		m.setStatus("copy failed: " + err.Error())
		return
	}
	m.setStatus("logs copied to clipboard")
}

func (m *Manager) save() {
	dump := m.buf.Dump()
	path, err := debugio.SaveDumpToDisk(dump)
	if err != nil {
		m.setStatus("save failed: " + err.Error())
		return
	}
	m.setStatus("logs saved to " + path)
}
