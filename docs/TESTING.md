hi hi gonna do a bit of a brainstorm here about what/how we'd go about testing things

## What needs to work?
- [ ] main page needs to load
- [ ] search queries need to be able to be made
  - [ ] and results need to be returned from those queries
- [ ] mhtml needs to be parsed/rendered back into html properly

## How do you propose we test these things?
that's a great question. The only test frameworks I've used in the past use an assert to make sure the result of the test is as expected.
Go has a test library built into the std. how does that library work?
Very similarly to the assert based ones. It's just an if check instead and then you call an error/failed method of the testing module when something went wrong.

### So how do we test the above things?
- [ ] main page needs to load
  - [ ] make an http request to the endpoint and check that a 200 code is written
    - this would work but it doesn't really check that the page is doing what it's supposed to
    - this also requires that the program be up and running already. Is there a way to do this easily?
- [ ] search queries need to be able to be made and results need to be returned from those queries
  - [ ] make an http request to the endpoint and inspect the response code and contents
- [ ] mhtml needs to be parsed/rendered back into html properly
  - [ ] we need a known working parser (cough chrome) and to test the mhtml dump against whatever chrome spits out
    - problems with this is there must be so so so many edge cases. I want to say we don't need to build an entire mhtml parser but we kinda might? this one might not be worth testing right now because we can kinda go off of visual checks with such a small corpus, and if things break we can note the bug and the document with which it occurred and go from there.


