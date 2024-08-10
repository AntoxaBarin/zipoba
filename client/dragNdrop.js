const SERVER_URL = 'http://127.0.0.1:8080/upload';

const dragArea = document.querySelector('.drag-area');
const dragText = document.querySelector('.header');

let button = document.querySelector('.button');
let input = document.querySelector('input');

let file;

button.onclick = () => {
    input.click();
}

input.addEventListener('change', function () {
    file = this.files[0];
    console.log('Uploading file...');
    uploadFile(file);
});

function uploadFile(file) {
    let formData = new FormData();
    formData.append('file', file);

    fetch(SERVER_URL, {
        method: 'POST',
        body: formData,
    }).then(response => {
        console.log(response.status);
        if (!response.ok) {
            alert('Network response was not ok');
        }
        return response.json();
    }).then (data => {
        console.log('Success: ', data);
    })
}

dragArea.addEventListener('dragover', (event) => {
    event.preventDefault();
    dragText.textContent = 'Release to upload';
    dragArea.classList.add('active');
});

dragArea.addEventListener('dragleave', () => {
    dragText.textContent = 'Drag-and-drop';
    dragArea.classList.remove('active');
});

dragArea.addEventListener('drop', (event) => {
    event.preventDefault();

    file = event.dataTransfer.files[0];
    console.log(file);

    let fileReader = new FileReader();
    fileReader.onload = () => {
        let fileURL = fileReader.result;
    }
    fileReader.readAsDataURL(file);
});
