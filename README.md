# GoSeq  
*A Note Tool for Developers*

GoSeq helps you manage daily and project-based notes, along with tracking TODOs in your code, similar to the `fixme` tag.

## Table of Contents
1. [Features](#features)
2. [Installation](#installation)  
   - [GitHub Credentials](#github-credentials)
   - [Adding Autocompletion](#adding-autocompletion)
3. [How It Works](#how-it-works)  
   - [Daily Notes](#1-daily-notes)  
     - [Create Daily Note](#create-daily-note)  
     - [Listing Notes](#listing-notes)  
     - [Join Notes](#join-notes)  
     - [Search Within Notes](#search-within-notes)  
     - [Delete Notes](#delete-notes)  
   - [Projects/Repos](#2-projectsrepos)  
     - [Project Notes](#project-notes)  
     - [List Projects](#list-projects)  
     - [Delete Projects](#delete-projects)  
   - [TODO Tracking](#3-todo-tracking-todooos)  
     - [Urgency System](#urgency-system)  
     - [Scan TODOs](#scan-todos)  
     - [Post TODOs](#post-todos)
4. [Help](#help)
5. [Lock Error](#lock-error)

---

## Features
- **Daily Notes Management**  
    - Create and manage notes for any time period.
- **Project-Based Notes**  
    - Keep notes linked to specific projects or repositories.
- **TODO Management**  
    - Tracks TODOs in your codebase, automatically pushing issues to GitHub.

---

## Installation  
Install GoSeq using the following command:

```bash 
go install github.com/DnFreddie/goseq@latest
```

### GitHub Credentials  
GoSeq retrieves your GitHub credentials from `$HOME/.config/.GoSeq.yaml` or `$HOME/.GoSeq`.  
If no credentials are found, GoSeq will prompt you to provide them.

Example `.GoSeq.yaml` config file:

```yaml
token: <personal-token>
```

To generate a Personal Access Token, visit [GitHub Settings](https://github.com/settings/tokens).  
Ensure the token has full access to private repositories.

---

### Adding Autocompletion

To enable command autocompletion:

```bash
# create the user completion directory
mkdir ~/.bash_completion.d/
# Generate the completion
goseq completion bash > ~/.bash_completion.d/goseq_completion.sh
# Make it executable
chmod +x ~/.bash_completion.d/goseq_completion.sh
```

Add the following lines to your `.bashrc`:

```bash
# Load custom bash completions
if [ -d ~/.bash_completion.d ]; then
    for file in ~/.bash_completion.d/*; do
        source "$file"
    done
fi
```

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
    - The default is one week.  
- Any changes made to the combined note are applied to the original notes.

```bash
goseq join -r <week|day|year|all> -t 3
```

#### Search Within Notes  
Search for specific patterns within your notes:

```bash
goseq search test agenda
```

##### Case Insensitive Search: `-i`

```bash
goseq search -i test agenda  
```

##### Regex Search: `-E`

```bash
goseq search -E ^test$agenda
```

##### Combining Flags

```bash
goseq search -i -E ^TEST$AGENDA
```

#### Delete Notes  
Open joined files in the editor and remove the one deleted in the document by the user.

```bash 
goseq delete 
```

---

### 2. Projects/Repos

#### Project Notes  
- Use the `-r` flag to open the most recently accessed project note.  
    - If no recent note is found, you’ll be prompted to choose a project manually:

```bash
goseq git -r 
```

- Or, provide a path to the directory containing the repository or project:

```bash
goseq git --path <Repo/dir containing repos>
```

This adds the project to the `$HOME/Documents/Agenda/projects/.PROJECTS_META.json` file.

#### List Projects  
To list added projects:

```bash 
goseq git list
```

#### Delete Projects  
Opens the names of the joined files inside the editor, removing those deleted in the document by the user.

```bash 
goseq git delete 
```

---

### 3. TODO Tracking (TODOOOS)

GoSeq finds and tracks TODOs in your project, compares them with existing TODOs, and pushes any new issues to GitHub.

#### Urgency System  
The urgency system is adapted from the [Fixmee Emacs extension](https://github.com/rolandwalker/fixmee#explanation).  
The urgency of a TODO is indicated by repeating the final character of the keyword (e.g., TODOOOO for a critical issue). The `scan` command sorts TODOs based on their urgency.

#### Scan TODOs  
Search for TODOs in the provided directory and generate a report:

```bash
goseq git scan -p <Repo/dir containing repos>
```

- Add the project to the known project list after scanning:

```bash
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

And push the ones that do not exist.

```bash
goseq git post  -p <path/to/the/repo>
```
![Goseq Post Example](/public/static/goseqPlan.png)


---

## Help 
For help with any commands, simply run `goseq` without arguments:

```bash
goseq
```

---

## Lock Error  
GoSeq uses lock files stored in `/tmp/`. If the program stops with the message `log not acquired`, you may need to delete the lock manually or reboot the system.

**Lock Files:**
- `/tmp/.goseq_delete.lock`
- `/tmp/.goseq_project_delete.lock`
- `/tmp/.goseq_join.lock`
- `/tmp/.goseq_project_join.lock`

--- 
