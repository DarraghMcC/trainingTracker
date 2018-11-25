# Calendar Training Counter #

Golang app to track Brazilian jiu jitsu training frequency. 
This app pulls from personal google "trained" calender all events labelled "BJJ" and graphs them for the user to see.

### Running the application ###

1. Clone the repository
2. Run ```./trainingTracker ```
3. View the output.png

### Current app output ###
![](https://github.com/DarraghMcC/trainingTracker/blob/master/output.png?raw=true)


### Things to note ###
* Credentials and tokens are not committed to this public repository. To build this project one must set up their own keys against their own calender.
* This is a personal project for learning GO and is very much a work in progress

### Known TODOs ###
* X-axis labelling is off
* Add some form of unit tests
* Paramterise which calenders and keys words are plotted against
* Add in plots for non Bjj training tasks (jogging + gym for example)
* Reduce console logging noise



