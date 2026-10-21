# uWoofer
okay so like- what if you cut out the webserver part and made this thing a native app? let's do a little microwoofer experiment.
well bubbletea exists. It'd be basically the same but instead of having to parse mhtml for display we might just be able to have the os open it in chrome. Can we do that? let's find out!

Okay so yes we totally can do that.

todos:
- [ ] make a way to store `config.toml` in `.metadata`
- [ ] ASYNC
  - [ ] do corpus launch n stuff asynchronously
    - [ ] launch tui
    - [ ] load config, index & db in a goroutine
    - [ ] connect them to the model once loaded
    - [ ] handle search coming in before background is loaded
      - [ ] either queue search or respond with an error notification
- [ ] UX
  - [ ] displaying blurbs
  - [ ] modifying documents
  - [ ] deleting documents
  - [ ] adding documents
  - [ ] opening file directly (remove extra scrolling time)
    - [ ] when in results context menu hitting `/` should open a text box. (this could become a commands section, but I'll leave that for later)
    - [ ] typing a number in this text box and hitting enter should:
      - [ ] try to set `m.cursorPos` to that value
        - [ ] if number is out of range, display out of range error, or ignore it
      - [ ] if `m.cursorPos` was successfully set, open the result at position in the results list.
  - [ ] colors/options configuration
    - [ ] this'd be controlled in `config.toml` but that's out of scope atm.
- [ ] KEEPING THE REPO CLEAN
  - [ ] fork uWoofer off of miniwoofer.




# Search UX
Make a search and then the list of titles will show up.
then you should be able to open a file by scrolling to it w/arrows or kj. pressing enter will open the selected file on your computer.

## Selecting An Element

# Keybinds / Concepts
## Search Mode
When you open uWoofer you're dumped into Search mode by default which puts your focus on the search bar. typing in this mode will type in the search bar and pressing enter will search whatever is in the search box. You know you're in search mode when your cursor is blinking in the search bar.
## Select Mode
Once you do a search you are dumped into select mode. This mode lets you scroll through results with the arrow keys or kj. pressing enter will open a selected file on your computer. to get back to search mode from here you can click `esc`. To get to select mode from search you can also use the up arrow


