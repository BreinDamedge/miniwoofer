# On Issue #35:
```txt
figure out:
- [ ] what needs testing
- [ ] what framework will be used
- [ ] how implementation should work
  - [ ] std testing vs 3rd party lib
- [ ] test corpus download script
```

That's what I wrote in the issue. I'd like to take a moment to brainstorm what needs testing in this application and how we'd go about doing that.

Since the core of it is a webserver we need a way to test all the endpoints and measure the actual behavior of the app vs the expected behavior. Another thing that needs testing is the mhtml -> html parser. Chrome can parse mhtml files already so one method could be to compare the output of our parser with the output of chrome's. Alternatively if there was a way to just have chrome do it and then display that, that'd be ideal. no, don't do that because then it needs more setup from the user, and it's something out of our control that would then need maintaining. So back to testing.

1. We need to test our endpoints
  - what are our endpoints?
    - [ ] display file                (GET /corpus/{file...})
    - [ ] text search                 (GET /)
    - [ ] rescan trigger              (GET /triggers/rescan)
    - [ ] open file management widget (GET /management.html)
    - [ ] delete a file               (POST /triggers/delete)
    - [ ] modify file blurb           (POST /triggers/update_blurb)
    - [ ] modify file title           (POST /triggers/update_title)

1. We need to test our parsing
  - what are we parsing?
    - [ ] html -> plain text (for indexing)
    - [ ] mhtml -> html (for rendering in browser)
