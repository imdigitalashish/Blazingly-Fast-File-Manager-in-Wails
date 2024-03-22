import './style.css';
import './app.css';

import logo from './assets/images/logo-universal.png';
import { Greet, GetAllDrives, GetAllFolders } from '../wailsjs/go/main/App';



class Application {
    constructor() {

        this.driveContainers = document.querySelector(".driveContainers");
        this.main();
        this.currentPath = [];

        this.directoryList = {};
    }

    utils = {
        getAllDrives: async () => {
            let request = await GetAllDrives();
            return JSON.parse(request);
        },
        getAllFoldersFromAPath: async (path) => {
            let request = await GetAllFolders(path);
            return JSON.parse(request);
        },
        changeDirectoryList: (path, data) => {
         
        }
    }

    main() {
        this.utils.getAllDrives()
            .then((res) => {
                Object.keys(res).forEach((val) => {
                    let driveBox = document.createElement("div");
                    driveBox.classList.add("driveBox")
                    let driveLetter = document.createElement("p")
                    driveLetter.innerText = `Local Disk (${val}:)`
                    driveBox.appendChild(driveLetter)


                    let progressBar = document.createElement("div")
                    progressBar.classList.add("progressBar");
                    let progressBarRate = document.createElement("div")
                    progressBarRate.classList.add("progressBarRate");

                    let width = ((res[val].total_space - res[val].space_left) / res[val].total_space)
                    if (width > 0.8) progressBarRate.style.backgroundColor = "red";
                    progressBarRate.style.width = (width * 200) + "px";
                    progressBar.appendChild(progressBarRate)
                    driveBox.appendChild(progressBar);
                    driveBox.dataset.drive = val;

                    driveBox.addEventListener("click", (e) => {
                        let drive = e.currentTarget.dataset.drive+":";
                        this.currentPath.push(drive);
                        console.log(this.currentPath.join("//"))
                        this.utils.getAllFoldersFromAPath(this.currentPath.join("//"))
                            .then((res) => {
                                console.log(res)
                            })
                    })

                    let para = document.createElement("p")
                    para.innerText = `${res[val].space_left} GB freed of ${res[val].total_space} GB`
                    driveBox.appendChild(para)

                    this.driveContainers.appendChild(driveBox)
                })




            })
    }

}



addEventListener("DOMContentLoaded", () => {
    window.app = new Application();
})



