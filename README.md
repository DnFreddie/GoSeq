# GoSeq  
*A Note Tool for Developers*

GoSeq helps you manage daily and project-based notes, along with tracking TODOs in your code, similar to the `fixme` tag.

## Features:
- **Daily Notes Management**  
    - Create and manage notes for any time period.
- **Project-Based Notes**  
    - Keep notes linked to specific projects or repositories.
- **TODO Management**  
    - Tracks TODOs in your codebase, automatically pushing issues to GitHub.
---

## How It Works

### 1. Daily Notes  
GoSeq creates a note for each day, which can be queried or combined as needed.

#### Create Daily Note  
- Open a new daily note or access an existing one.  
- Notes are stored in `$HOME/Documents/Agenda/`.
```bash
goseq new
```

#### Listing Notes  
- List all daily notes and choose the one you wish to edit:
```bash
goseq list
```

#### Join Notes  
- Combine multiple notes (from a specific period like a week or a year) into one.  
    -  the defualt is from one week 
- Any changes made to the combined note are applied to the original notes.

`-r` is for datetime
`-t` How many times ex. notes from 3 weeks

```bash
goseq join -r <week|day|year|all> -t 3
```

---

#### Search Within Notes

This feature allows you to search for specific patterns within your notes and select the one you wish to open.

```bash
# GoSeq will join the results, so you don't need to worry about the quotes.
goseq search test agenda
```

The flags used are similar to those in `grep`, making them familiar to users.

##### Case Insensitive Search: `-i`

```bash
goseq search -i test agenda  
```

#####  Regex: `-E`

```bash
goseq search -E ^test$agenda
```

---

##### Combining Flags

You can combine both flags for more flexible searching:

```bash
goseq search -i -E ^TEST$AGENDA
```


#### Delete Notes

It opens a names  of the joined files inside the editor.
And removes the one that has been deleted in the document by the user.

```bash 
goseq delte 
```

--- 

### 2. Projects/Repos

#### Project Notes  
-  `-r` Open the most recently accessed project note.  
    - If no recent note is found, you’ll be prompted to choose a project manually:

```bash
#Open a recent Project
goseq git -r 
```
- Optionally, provide a path to the directory containing the repository or project:

```bash
goseq git --path <Repo/dir containing repos>
```
`git --path` adds the project to the file called  `$HOME/Documents/Agenda/projects/.PROJECTS_META.json`

**List Projects**
- To list porject that were added use 

```bash 
goseq git list
```

---

### 3. TODO Tracking (TODOOOS)  

GoSeq finds and tracks TODOs in your project, compares them with existing TODOs, and pushes any new issues to GitHub.

#### Urgency System  
The urgency system is adapted from the [Fixmee Emacs extension](https://github.com/rolandwalker/fixmee#explanation).  
The urgency of a TODO is indicated by repeating the final character of the keyword (e.g., TODOOOO for a critical issue). The `scan` command sorts TODOs based on their urgency.

#### Scan TODOs  
- Search for TODOs in the provided directory and generate a report:
```bash
goseq git scan -p <Repo/dir containing repos>
```
- `-a` add the project to the Known project after scaning

```bash
#if dir has more repos it will also save them
goseq git scan -a -p <Repo/dir containing repos>
```

##### Example Report 
```md
Project: DnFreddie/Blog
------------------------------
Location: drive.svelte  
TODO: Fix the animation on the banner  
Line: 2  
Urgency: 5  
------------------------------
```

#### Post TODOs  
Post new TODO issues to GitHub
This will check weather todos already exist on the github.
Then ask you do you want to push them.
And push the onese that do not exist.

```bash
goseq git post  -p <path/to/the/repo>
```
![Goseq Post Example](/public/static/goseqPlan.png)

---


### GitHub Credentials  
GoSeq retrieves your GitHub credentials from `$HOME/.config/.GoSeq.yaml` or `$HOME/.GoSeq`.  
If no credentials are found, GoSeq will prompt you to provide them. Example config:

```yaml
token: <personal-token>
```

To generate a Personal Access Token, visit [GitHub Settings](https://github.com/settings/tokens).  
Ensure the token has full access to private repositories.

---

## Usage  
For help with any commands, simply run `goseq` without arguments:

```bash
goseq
```

---

## Build

To build GoSeq:

```bash
export GOPATH=$PWD
go get .
go build .
```


