// Package tui provides a tui component for uWoofer.
package tui

import (
	"fmt"
	"math/rand"
	"strings"

	// for tui
	"charm.land/bubbles/v2/cursor" // search box prettyification
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2" // general tui
	"charm.land/lipgloss/v2"

	// for search
	"mime"
	"miniwoofer/impl"
	"regexp"

	// for opening file
	"os/exec"
	"runtime"
)

func openDocument(filePath_ string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// Windows
		cmd = exec.Command("cmd", "/c", "start", filePath_)
	case "darwin":
		// macOS
		cmd = exec.Command("open", filePath_)
	default:
		// Linux
		cmd = exec.Command("xdg-open", filePath_)
	}

	return cmd.Start()
}

type searchResult struct {
	score float64
	meta  *impl.DocumentMeta
}

type Model struct {
	index         *impl.Bm25
	metaDb        *impl.MetaDb
	config        impl.Config
	cursorPos     int
	searchResults []searchResult
	searchBox     textarea.Model
	viewport      viewport.Model
	scrollDist    int
	lastSearch    string // string version of last search performed
	senderStyle   lipgloss.Style
	err           error
}

func initializeSearchComponents() (impl.Config, *impl.Bm25, *impl.MetaDb, error) {
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
	return config, index, db, nil
}

func InitialModel() Model {
	ta := textarea.New()
	ta.Placeholder = "Search..."
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
	vp.SetContent(`uWoofer baybee!!!`)
	vp.KeyMap.Left.SetEnabled(false)
	vp.KeyMap.Right.SetEnabled(false)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	config, index, db, err := initializeSearchComponents()

	return Model{
		// all the things the tui needs
		searchResults: []searchResult{},
		searchBox:     ta,
		lastSearch:    "",
		cursorPos:     -1, // when no searches exist
		viewport:      vp,
		scrollDist:    3,
		senderStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		index:         index,
		metaDb:        db,
		config:        config,
		err:           err,
	}
}

func random_sentence(length_ int) string {
	vocab := strings.Split("the quick brown fox is getting pretty lazy so he doesn't want to jump over any more logs but maybe they will eat some soup because she's kinda tired tbh and soup sounds nice to a tired slow lil guy", " ")
	var output strings.Builder
	for range length_ {
		randomIndex := rand.Intn(len(vocab))
		output.WriteString(vocab[randomIndex])
		output.WriteString(" ")
	}
	return output.String()
}

func (m *Model) doSearch() {
	// get query
	var query string = m.searchBox.Value()

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

	results, err := m.index.Search(search_terms) // returns results already ranked
	if err != nil {
		panic(err)
	}

	// perform search
	// populate the model's search_results field
	m.searchResults = []searchResult{} // reset it first
	for _, s := range results {
		var sr searchResult
		sr.score = s.Score

		// get title from metadata
		doc_md, err := m.metaDb.GetDocument(s.Id)
		if err != nil {
			panic(err)
		}
		sr.meta = doc_md
		// now put it in the model for state
		m.searchResults = append(m.searchResults, sr)
	}

	// set cursorPos to top
	m.cursorPos = 0
}

func (m Model) renderSearchResults() string {
	if len(m.searchResults) == 0 {
		return "No results!"
	}

	var br strings.Builder
	// br.WriteString("Results:\n")
	for i, result := range m.searchResults {
		cursor := " "
		if m.cursorPos == i {
			cursor = ">"
		}
		// br.WriteString(fmt.Sprintf("%s %d. %s\n", cursor, i+1, result)) // this is slower than the following (thank you linter)
		fmt.Fprintf(&br, "%s %d. %s\n", cursor, i+1, result.meta.Title)
	}
	return br.String()
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// focus command catching
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		// scrolling results
		if !m.searchBox.Focused() {
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "up", "k":
				if m.cursorPos > 0 {
					m.cursorPos--
					// scroll viewport and re-render new text
					m.viewport.ScrollUp(1)
					m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width()).Render(m.renderSearchResults())) // TODO: consolidate this
				}
			case "down", "j":
				if m.cursorPos < len(m.searchResults)-1 {
					m.cursorPos++
					// scroll viewport and re-render new text
					m.viewport.ScrollDown(1)
					m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width()).Render(m.renderSearchResults())) // TODO: consolidate this
				}
			case "enter":
				// open result under cursor
				// get file id
				file_id := m.searchResults[m.cursorPos].meta.Id
				// use open command
				openDocument(m.config.CorpusDir + "/" + file_id) // TODO: string sanitization (mostly 4 pathsep)
			case "esc":
				// focus the search bar
				m.searchBox.Focus()
			}
		} else {
			// type into text box...
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "up":
				// unfocus and return to searchResults
				m.searchBox.Blur()
			case "enter":
				// send search
				m.doSearch()
				m.searchBox.Blur()
				m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width()).Render(m.renderSearchResults())) // TODO: consolidate this
			default:
				// Send all other keypresses to the searchBox textarea.
				var cmd tea.Cmd
				m.searchBox, cmd = m.searchBox.Update(msg)
				return m, cmd
			}
		}

	case tea.WindowSizeMsg:
		m.viewport.SetWidth(msg.Width)
		m.searchBox.SetWidth(msg.Width)
		m.viewport.SetHeight(msg.Height - m.searchBox.Height())

		if len(m.searchResults) > 0 {
			// Wrap content before setting it.
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width()).Render(m.renderSearchResults())) // render is where your string layout goes
		}
		m.viewport.GotoTop()

	case cursor.BlinkMsg:
		// searchBox textarea should also process cursor blinks.
		var cmd tea.Cmd
		m.searchBox, cmd = m.searchBox.Update(msg)
		return m, cmd
	}

	// cleanup (handles refocus)
	m.searchBox, cmd = m.searchBox.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() tea.View {
	viewportView := m.viewport.View()
	v := tea.NewView(viewportView + "\n" + m.searchBox.View())
	c := m.searchBox.Cursor()
	if c != nil {
		c.Y += lipgloss.Height(viewportView)
	}
	v.Cursor = c
	v.AltScreen = true
	return v
}
