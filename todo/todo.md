1. Code is under documented we need to quickly fix this. which may be done as follows.
2. use script to stream and break go files into 4000 chunks make gpt api call to reduce/simpify block (replace repetive sections  in lists with [..] , remove new lines, etc. while mainting origina line number history)
3. concat slimmed down files into staging dir, then proceed to break up and stream into gpt4 model with prompts instructing it review and summarize the individual files, whle listing public members, methods, classes and method descriptions, concerns, suggestions into a yaml format (and some errata) 
4. From this second generation data tweak feead back and forth into gpt4 briefly to get a full project overview and per file/section overview. feeding it back in the previously intermediate yaml to assist it in constructing those details.
5. Update codebase. 
6. Feed in the very high level project summary to GPT4 + per section detailed yaml with a prompt to provide todo/follow up items/improvements.   
etc. 
7. Feed in raw file segments again with high level and midlevel summarizing yaml to construct/generate inline commets/notes.

I have some fledgling tools and prompts to semi automate this process I'll provide later. 