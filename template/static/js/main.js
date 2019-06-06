const container = document.querySelector("#container");

function createItem(id) {
    const newElement = document.createElement('div');
    newElement.setAttribute('id', 'moving-div-' + id);
    newElement.setAttribute('class', 'movable');
    container.appendChild(newElement);
    newElement.addEventListener('click', function(e) {
        pickingItem(e.target.id.replace('moving-div-',''));
    })
}

function removeItem(id) {
    for(let i = 0; i <  container.childElementCount; i++) {
       if (container.children[i].id === 'moving-div-' + id) {
           container.removeChild(container.children[i]);
           return;
       }
    }
}

async function pickingItem(id) {
    const data = await fetch('http://localhost:8080/picking/' + id);
    const responseBody = await data.json();
    console.log(responseBody);
    removeItem(id);
}

async function initializeItemList() {
    try {
        const data = await fetch('http://localhost:8080/items');
        const items = await data.json();
        items.forEach(item => createItem(item.Id));
    } catch (e) {
        console.error(e);
    }

}

initializeItemList();