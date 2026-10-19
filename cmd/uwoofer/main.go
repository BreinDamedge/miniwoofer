package main

// main file for μWoofer experiment

import (
	// general os stuff
	"fmt"
	"os"
	"strings"

	// for bubble tea
	"charm.land/bubbles/v2/cursor"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	// for search
	"mime"
	"miniwoofer/impl"
	"regexp"
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Oof: %v\n", err)
	}
}

type search_result struct {
	score float64
	meta  *impl.DocumentMeta
}

type model struct {
	viewport       viewport.Model
	search_results []search_result
	textarea       textarea.Model
	senderStyle    lipgloss.Style
	b              *impl.Bm25
	db             *impl.MetaDb
	err            error
}

func initialModel() model {
	index := impl.NewBm25()
	config := impl.LoadConfig()

	// add mime extentions that may not exist by default
	mime.AddExtensionType(".mht", "multipart/related")
	mime.AddExtensionType(".mhtml", "multipart/related")

	db, err := impl.MetaDbOpen(config)
	if err != nil {
		panic(err)
	}

	fmt.Print("Attempt to load serialized BM25... ")
	if err := index.Load(config.Bm25Path()); err != nil {
		fmt.Println("Serialized BM25 Does not exist.")
		if err := impl.ParseCorpus(index, config); err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Success.")
	}

	fmt.Print("Checking if corpus changed... ")
	if changed, err := impl.CheckChanged(config.CorpusDir, index.CorpusHash); err != nil || changed {
		fmt.Print("Corpus Changed\n")
		if err := impl.ParseCorpus(index, config); err != nil {
			panic(err)
		}
		if err := db.AddCorpus(config); err != nil {
			panic(err)
		}
	} else {
		fmt.Print("No Change.\n")
	}

	ta := textarea.New()
	ta.Placeholder = "query...?"
	ta.SetVirtualCursor(false)
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(1)

	// Remove cursor line styling
	s := ta.Styles()
	s.Focused.CursorLine = lipgloss.NewStyle()
	ta.SetStyles(s)

	ta.ShowLineNumbers = false

	vp := viewport.New(viewport.WithWidth(30), viewport.WithHeight(5))
	vp.SetContent(`Welcome to the μWoofer experiment.`)
	vp.KeyMap.Left.SetEnabled(false)
	vp.KeyMap.Right.SetEnabled(false)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		textarea:       ta,
		search_results: []search_result{},
		viewport:       vp,
		senderStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		b:              index,
		db:             db,
		err:            nil,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.SetWidth(msg.Width)
		m.textarea.SetWidth(msg.Width)
		m.viewport.SetHeight(msg.Height - m.textarea.Height())

		if len(m.search_results) > 0 {
			// what if we do nothing here?
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width()).Render(""))
		}
		m.viewport.GotoBottom()
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case "enter": // here's ur search query
			// get query
			var query string = m.textarea.Value()

			// maybe do a check to make sure it's not empty...

			// search and populate search results

			// stemming ig?
			re := regexp.MustCompile(`((?:".*?" *)|(?:(?:\S+? +)))`)
			search_terms := []string{}
			for _, match := range re.FindAllStringSubmatch(query+" ", -1) {
				if len(match) > 1 {
					search_terms = append(search_terms, strings.Trim(strings.ToLower(match[1]), ` "`))
				}
			}

			results, err := m.b.Search(search_terms) // returns results already ranked
			if err != nil {
				panic(err)
			}

			// render the results
			if len(results) == 0 {
				m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width()).Render("No Results for \"" + m.textarea.Value() + "\""))
			} else {
				// there are search results so you should make that list of titles and render them.
				// titles can be retrieved from the metadata db given the ids of the results

				// perform search
				// populate the model's search_results field
				// and build displayed string in the same loop
				m.search_results = []search_result{} // reset it first
				var br strings.Builder
				for i, s := range results {
					var sr search_result
					sr.score = s.Score

					// get title from metadata
					doc_md, err := m.db.GetDocument(s.Id)
					if err != nil {
						panic(err)
					}
					sr.meta = doc_md
					// now put it in the model for state
					m.search_results = append(m.search_results, sr)
					br.WriteString(fmt.Sprintf("%d. %s\n", i, sr.meta.Title))
				}

				// display the search results
				m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width()).Render(br.String()))
			}
			m.textarea.Reset()
			m.viewport.GotoBottom()
			return m, nil
		default:
			// Send all other keypresses to the textarea.
			var cmd tea.Cmd
			m.textarea, cmd = m.textarea.Update(msg)
			return m, cmd
		}

	case cursor.BlinkMsg:
		// Textarea should also process cursor blinks.
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() tea.View {
	viewportView := m.viewport.View()
	v := tea.NewView(viewportView + "\n" + m.textarea.View())
	c := m.textarea.Cursor()
	if c != nil {
		c.Y += lipgloss.Height(viewportView)
	}
	v.Cursor = c
	v.AltScreen = true
	return v
}
