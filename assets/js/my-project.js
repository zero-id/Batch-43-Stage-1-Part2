let project = [];

//get data

function getData(event) {
  event.preventDefault();

  let projectName = document.getElementById("project-name").value;
  let startDate = document.getElementById("start-date").value;
  let endDate = document.getElementById("end-date").value;
  let description = document.getElementById("description").value;
  let image = document.getElementById("img").files;

  console.log(image);

  image = URL.createObjectURL(image[0]);

  // technologiesChecked

  let nodeJs = document.getElementById("node-js");
  let reactJs = document.getElementById("react-js");
  let nextJs = document.getElementById("next-js");
  let typeScript = document.getElementById("type-script");

  let iconNodeJs = ``;
  let iconReactJs = ``;
  let iconNextJs = ``;
  let iconTypeScript = ``;

  if (nodeJs.checked === false) {
    iconNodeJs = `display: none;`;
  }
  if (reactJs.checked === false) {
    iconReactJs = `display: none;`;
  }
  if (nextJs.checked === false) {
    iconNextJs = `display: none;`;
  }
  if (typeScript.checked === false) {
    iconTypeScript = `display: none;`;
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
    iconTypeScript,
  };

  project.push(addProject);
  console.log(project);
  showData();
}

//shoeData

function showData() {
  document.getElementById("containerList").innerHTML = ` `;

  for (let i = 0; i < project.length; i++) {
    document.getElementById("containerList").innerHTML += `
    <div class="w-25 h-75 d-flex flex-column p-3 gap-3 flex-nowarp shadow-lg rounded-3">
            <div class="">
              <img class="w-100 h-100 rounded-3" src="${project[i].image}" alt="" />
            </div>
            <div class="">
              <h5><a href="project-detail.html" class="text-black fw-bold">${project[i].projectName}</a></h5>
              <p class="fs-6 text-secondary">${duration(project[i].startDate, project[i].endDate)}</p>
            </div>
            <div style="height: 100px" class="overflow-hidden">
              <p class="">${project[i].description}</p>
            </div>

            <div class="d-flex gap-2 mb-2">
              <img style="width: 40px; ${project[i].iconNextJs}" class="" src="assets/img/nextJs.png" />
              <img style="width: 40px; ${project[i].iconNodeJs}" class="" src="assets/img/nodeJs.png" />
              <img style="width: 40px; ${project[i].iconReactJs}" class="" src="assets/img/reactJs.png" />
              <img style="width: 40px; ${project[i].iconTypeScript}" class="" src="assets/img/typeScript.png" />
            </div>
            <div class="d-flex container gap-2">
              <button class="btn btn-dark w-50">edit</button>
              <button class="btn btn-dark w-50">delet</button>
            </div>
          </div>`;
  }
}

//time duration

function duration(startDate, endDate) {
  let selisih = new Date(endDate) - new Date(startDate);

  let day = Math.floor(selisih / (1000 * 60 * 60 * 24));
  let month = Math.floor(selisih / (1000 * 60 * 60 * 24 * 30));

  if (month > 0) {
    return `Durasi: ${month} Bulan`;
  } else {
    return `Durasi: ${day} Hari`;
  }
}
