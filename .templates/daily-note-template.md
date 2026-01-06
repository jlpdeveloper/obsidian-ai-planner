---
tags:
  - daily-log
rating: 
icon: RiQuillPenLine
---
> [!multi-column]
> 
> >[!todo] Weekly Goals
> >```dataview
 >>task
> >from #weekly-log
> >where journal-date <= this.journal-date and journal-end-date >= this.journal-date
> >```
 > 
 > > [!summary] Yesterday's Log
 > > ```dataview
>>TABLE WITHOUT ID file.link AS "Link",
>>""+length(filter(file.tasks, (r) => contains(list(["x", "!", "f", "<", "k", "i"]), r.status) and meta(r.header).subpath = "Goals" and length(r.children) = 0))+"/" + length(filter(file.tasks, (r) => !contains(list(["-", "?", "/"]), r.status) and meta(r.header).subpath = "Goals" and length(r.children) = 0)) AS Planned,
>>""+length(filter(file.tasks, (r) => contains(list(["x", "!", "f", "<", "k", "i"]), r.status) and meta(r.header).subpath = "Meetings" and length(r.children) = 0))+"/"+length(filter(file.tasks, (r) => !contains(list(["-", "?", "/"]), r.status) and meta(r.header).subpath = "Meetings" and length(r.children) = 0)) AS Meetings,
>>""+length(filter(file.tasks, (r) => contains(list(["x", "!", "f", "<", "k", "i"]), r.status) and meta(r.header).subpath = "Bonus Items" and length(r.children) = 0))+"/"+length(filter(file.tasks, (r) =>  !contains(list(["-", "?", "/"]), r.status) and meta(r.header).subpath = "Bonus Items" and length(r.children) = 0)) AS Unplanned,
>>length(filter(file.tasks, (r) => !contains(list(["?", "/", "-"]), r.status) and contains(list(["Goals", "Meetings", "Bonus Items"]), meta(r.header).subpath) and length(r.children) = 0)) AS "Total Items"
>>FROM #daily-log and -"templates"
>>WHERE date(file.name) < date(this.file.name)
>>SORT file.name DESC
>>LIMIT 1
>> ```

## Tasks
```dataview
TABLE WITHOUT ID 
"<progress max='100' value='" + round(100*length(filter(rows.tasks, (r) => contains(list(["x", "!", "f", "<", "k", "i"]), r.status)))/length(rows)) + "'>" + round(100*length(filter(rows.tasks, (r) => r.completed))/length(rows)) + "%</progress>" as "Progress Bar",
""+length(filter(rows.tasks, (r) => contains(list(["x", "!", "f", "<", "k", "i"]), r.status) and meta(r.header).subpath = "Goals"))+"/" + length(filter(rows.tasks, (r) => meta(r.header).subpath = "Goals")) AS Planned,
""+length(filter(rows.tasks, (r) => contains(list(["x", "!", "f", "<", "k", "i"]), r.status) and meta(r.header).subpath = "Meetings"))+"/"+length(filter(rows.tasks, (r) => meta(r.header).subpath = "Meetings")) AS Meetings,
""+length(filter(rows.tasks, (r) => contains(list(["x", "!", "f", "<", "k", "i"]), r.status) and meta(r.header).subpath = "Bonus Items"))+"/"+length(filter(rows.tasks, (r) => meta(r.header).subpath = "Bonus Items")) AS Unplanned,
length(rows) AS Total
WHERE file.path = this.file.path
FLATTEN filter(file.tasks, (r) => contains(list(["Goals", "Meetings", "Bonus Items"]), meta(r.header).subpath) and !contains(list([">","?", "-"]), r.status) and length(r.children) = 0) as tasks
Group by tasks.task
SORT length(rows) DESC
```
### Goals

### Meetings

### Bonus Items

## Notes

## Other
---
> [!multi-column] 
>> [!example]- Keyboard Shortcuts
>>![[Keyboard Shortcuts]]
>
>>[!example]- Links
>>![[Daily Note Important Links]]

![[DailyNotesViews.base]]
