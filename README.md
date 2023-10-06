### Apahe Scraper

A somewhat fun project of me scraping animepahe with playwright-go.

### Installation
If you want to compile this, the following requirements are needed:
- Go
- ....

Then we clone this repo
```bash
git clone https://github.com/Clashkid155/Apahe_Scraper.git

cd Apahe_Scraper
# Build the application
go build -o apahe
# Run the program 
./apahe -h
```
### Or 
Download the already compiled one for your os. Check the release.

___

## Usage

```bash
./apahe -n "natruo" -q 1080p
```

Where 
- -n: means anime name.
- -q: means quality.

Use `apahe -h` to see the help. 


### TODO
* [ ] Implement a range function.
* [ ] Proper method of saving anime details.
* [ ] ....