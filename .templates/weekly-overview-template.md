---
icon: RiCalendarEventLine
tags:
  - weekly-log
---
## Week Rating
```dataview 
TABLE WITHOUT ID 
mood AS "Average Mood"
from #daily-log 
WHERE file.ctime >= this.journal-date and file.ctime <= this.journal-end-date and file.name != this.file.name 
GROUP BY ""
FLATTEN round(average(nonnull(rows.rating)),2) as mood
```
## Weekly Goals


## Fires
```dataview
TABLE WITHOUT ID file.link AS "Title", file.ctime as "Created On", file.mtime AS "Last modified at",  fire as "Fire Item"
WHERE file.name != this.file.name and file.folder != "templates" and file.ctime >= this.journal-date and file.ctime <= this.journal-end-date and fire
SORT file.mtime DESC
LIMIT 25
```
```dataview
TABLE WITHOUT ID file.link AS "Title", file.ctime as "Created On", file.mtime AS "Last modified at",  tasks.text as "Fire Item"
FROM #daily-log and -"templates"
WHERE date(file.name) >= this.journal-date and date(file.name) <= this.journal-end-date and file.tasks
FLATTEN file.tasks as tasks
WHERE contains(tasks.text, "fire::")
SORT file.mtime DESC
LIMIT 25
```
## Architectural Work

```dataview
TABLE WITHOUT ID file.link AS "Title", file.ctime as "Created On", file.mtime AS "Last modified at",  topics as "Topics"
FROM -#daily-log
WHERE file.name != this.file.name and file.folder != "templates" and file.ctime >= this.journal-date and file.ctime <= this.journal-end-date and role = "architect"
SORT file.mtime DESC
LIMIT 25
```
```dataview
TABLE WITHOUT ID architect as "Comment", file.link AS "Title", file.ctime as "Created On"
WHERE file.name != this.file.name and file.folder != "templates" and file.ctime >= this.journal-date and file.ctime <= this.journal-end-date and architect
SORT file.mtime DESC
LIMIT 25
```


```dataview
TABLE WITHOUT ID file.link AS "Title", status as "Status", followup-date as "Follow Up Date", file.ctime as "Created On", file.mtime AS "Last modified at",  topics as "Topics"
FROM #followup 
WHERE file.folder != "templates" and role = "architect" and followup-date
SORT followup-date DESC
LIMIT 25
```

## Reviews
```dataviewjs
let count = 0;
const pages = dv.pages("#daily-log")
  .where(p => !p.file.name.toLowerCase().includes("template"))
  .where(p => p["journal-date"] >= dv.current()['journal-date'])
  .where(p => p["journal-date"] <= dv.current()['journal-end-date']);

for (let page of pages) {
  for (let task of page.file.tasks) {
    // Check for completed, and presence of inline review property
    let match = task.text.match(/\(review::[^\)]+\)/i);
    if (task.completed && match) {
      count++;
    }
  }
}

dv.paragraph(`Completed reviews: **${count}**`);
```

> [!NOTE]- Review List
>```dataview
>task FROM #daily-log and -"templates" WHERE journal-date >= this.journal-date and journal-date <= this.journal-end-date and review and completed SORT file.ctime DESC LIMIT 25
>```

## Weekly Summary
### Overall
```dataview
TABLE WITHOUT ID 
length(filter(rows.tasks, (r) => meta(r.header).subpath = "Goals" and !contains(list([">", "-", " "]), r.status))) AS Planned,
length(filter(rows.tasks, (r) => meta(r.header).subpath = "Meetings" and !contains(list([">", "-", " "]), r.status))) AS Meetings,
length(filter(rows.tasks, (r) => meta(r.header).subpath = "Bonus Items" and !contains(list([">", "-", " "]), r.status))) AS Unplanned,
length(filter(rows.tasks, (r) => r.status = ">")) AS Moved,
length(filter(rows.tasks, (r) => r.status = "-")) AS Canceled,
length(filter(rows.tasks, (r) => !contains(list([">", "-", " "]), r.status))) AS Total
FROM #daily-log 
WHERE journal-date >= this.journal-date and journal-date <= this.journal-end-date
FLATTEN filter(file.tasks, (r) => contains(list(["Goals", "Meetings", "Bonus Items"]), meta(r.header).subpath) and !contains(list(["?"]), r.status) and length(r.children) = 0) as tasks
Group by tasks.task
SORT length(rows) DESC
```

### By Day
```dataview
TABLE WITHOUT ID dateformat(date(file.name), "cccc") as "Day of Week", file.link AS "Title", 
""+length(filter(file.tasks, (r) => contains(list(["x", "!", "f", "<", "k", "i"]), r.status) and meta(r.header).subpath = "Goals" and length(r.children) = 0))+"/" + length(filter(file.tasks, (r) => meta(r.header).subpath = "Goals" and !contains(list([">","?", "/", "-"]), r.status) and length(r.children) = 0)) AS Planned,
""+length(filter(file.tasks, (r) => contains(list(["x", "!", "f", "<", "k", "i"]), r.status) and meta(r.header).subpath = "Meetings" and length(r.children) = 0))+"/"+length(filter(file.tasks, (r) => meta(r.header).subpath = "Meetings" and !contains(list([">","?", "/", "-"]), r.status) and length(r.children) = 0)) AS Meetings,
""+length(filter(file.tasks, (r) => contains(list(["x", "!", "f", "<", "k", "i"]), r.status) and meta(r.header).subpath = "Bonus Items" and length(r.children) = 0))+"/"+length(filter(file.tasks, (r) => meta(r.header).subpath = "Bonus Items" and !contains(list([">","?", "-"]), r.status) and length(r.children) = 0)) AS Unplanned,
length(filter(file.tasks, (r) => contains(list(["Goals", "Meetings", "Bonus Items"]), meta(r.header).subpath) and !contains(list([">","?", "-"]), r.status) and length(r.children) = 0)) AS Total
FROM #daily-log 
WHERE journal-date >= this.journal-date and journal-date <= this.journal-end-date
SORT journal-date
LIMIT 10
```

```dataviewjs
var this_page = dv.current()
var pages = dv.pages('#daily-log').where(t => {
let d = dv.date(t.file.name);
return d >= this_page['journal-date'] && d <= this_page['journal-end-date'];
})
let badStatus = [">", "-"];
let completeStatus = ["x", "f",];
//dv.el("b", pages.length)
let totalTaskData = [0,0,0];
let completeTaskData = [0,0,0];
pages.forEach((page) => {

totalTaskData[0] += page.file.tasks.where(c => c.section.subpath == "Goals" && badStatus.indexOf(c.status) <0 && c.children.length == 0).length;
totalTaskData[1] += page.file.tasks.where(c => c.section.subpath == "Meetings" && badStatus.indexOf(c.status) <0 && c.children.length == 0).length;
totalTaskData[2] += page.file.tasks.where(c => c.section.subpath == "Bonus Items" && badStatus.indexOf(c.status) <0 && c.children.length == 0).length;
completeTaskData[0] += page.file.tasks.where(c => c.section.subpath == "Goals" && completeStatus.indexOf(c.status) >=0 && c.children.length == 0).length;
completeTaskData[1] += page.file.tasks.where(c => c.section.subpath == "Meetings" && completeStatus.indexOf(c.status) >=0 && c.children.length == 0).length;
completeTaskData[2] += page.file.tasks.where(c => c.section.subpath == "Bonus Items" && completeStatus.indexOf(c.status) >=0 && c.children.length == 0).length;

});

const chartData = {
    type: 'radar',
    data: {
        labels: ['Planned', 'Meetings', 'Unplanned'],
        datasets: [{
            data: totalTaskData,
            label: 'Total Weekly Tasks',
			borderColor: 'rgb(54, 162, 235)',
			pointBackgroundColor: 'rgb(54, 162, 235)',
            fill: true
            
        },{
            data: completeTaskData,
            label: 'Completed Weekly Tasks',
			backgroundColor: 'rgba(255, 99, 132, 0.2)',
		    borderColor: 'rgb(255, 99, 132)',
			pointBackgroundColor: 'rgb(255, 99, 132)',
			pointBorderColor: '#fff',
            fill: true
            
        }]
    },
    backgroundColor: 'rgba(256, 256, 0, 0.1)',
    options: {
	    
        scales: {
            r: {
                suggestedMin: 0,
                suggestedMax: 20,
                angleLines: {
	                color: 'white'
                },
                grid:{
	                color: 'white'
                }
            }
        }
    }
}
//dv.el("b", JSON.stringify(pages[0].file.tasks, null, 4))
window.renderChart(chartData, this.container);

```
