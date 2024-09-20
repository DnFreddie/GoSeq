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
- Any changes made to the combined note are applied to the original notes.

```bash
goseq join -r <week|day|year|all>
```

---

### 2. Projects/Repos

#### Project Notes  
- Open the most recently accessed project note.  
- If no recent note is found, you’ll be prompted to choose a project manually:
```bash
goseq git
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
goseq git scan --path <Repo/dir containing repos>
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

#### Push TODOs  
- Push new TODO issues to GitHub:
```bash
goseq git push
```

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


