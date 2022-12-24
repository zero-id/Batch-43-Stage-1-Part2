    
    
let project = []

//get data

function getData(event) {
    event.preventDefault()

    let projectName = document.getElementById("project-name").value
    let startDate = document.getElementById("start-date").value
    let endDate = document.getElementById("end-date").value
    let description = document.getElementById("description").value
    let image = document.getElementById("img").files

    image = URL.createObjectURL(image[0])

    // technologiesChecked

    let nodeJs = document.getElementById("node-js")
    let reactJs = document.getElementById("react-js")
    let nextJs = document.getElementById("next-js")
    let typeScript = document.getElementById("type-script")

    let iconNodeJs = ``
    let iconReactJs = ``
    let iconNextJs = ``
    let iconTypeScript = ``

    if (nodeJs.checked === false) {
        iconNodeJs = `style="display: none;"`
    }
    if (reactJs.checked === false) {
        iconReactJs = `style="display: none;"`
    }
    if (nextJs.checked === false) {
        iconNextJs = `style="display: none;"`
    }
    if (typeScript.checked === false) {
        iconTypeScript = `style="display: none;"`
    }
    
    let addProject = {
        projectName,
        startDate,
        endDate,
        description,
        image,
        iconNodeJs,
        iconReactJs,
        iconNextJs,
        iconTypeScript
    }

    project.push(addProject)
    console.log(project)
    showData()
}

//shoeData

function showData() {
    document.getElementById("containerList").innerHTML = ``
    
    for (let i=0; i < project.length; i++) {
        document.getElementById("containerList").innerHTML += `<div class="miniProject">
        <div class="img">
            <img src="${project[i].image}" alt="">
        </div>
        <div class="title-project">
           <h3><a href="project-detail.html">${project[i].projectName}</a> </h3>
            <p>${duration(project[i].startDate, project[i].endDate)}</p> 
        </div>
        <div class="content"> 
            <p>${project[i].description}</p>
        </div>
        
        <div id="icon">
            <img src="assets/img/nextJs.png"${project[i].iconNextJs} >
            <img src="assets/img/nodeJs.png" ${project[i].iconNodeJs}>
            <img src="assets/img/reactJs.png"${project[i].iconReactJs} >
            <img src="assets/img/typeScript.png"${project[i].iconTypeScript} >
        </div>
        <div class="btn">
            <button class="edit">edit</button>
            <button class="delet">delet</button>
        </div>

    </div>`
    }
}

//time duration

function duration(startDate, endDate) {
    let selisih = new Date(endDate) - new Date(startDate)
    
    let day = Math.floor(selisih / (1000 * 60 * 60 * 24))
    let month = Math.floor(selisih / (1000 * 60 * 60 * 24 * 30))
    
    if (month > 0) {
        return `Durasi: ${month} Bulan`
    } else {
        return `Durasi: ${day} Hari`
    }
}




