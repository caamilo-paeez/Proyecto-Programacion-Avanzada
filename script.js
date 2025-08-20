
// -------------------- Tabs --------------------
document.querySelectorAll('.tab').forEach(tab => {
    tab.addEventListener('click', () => {
        document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
        document.querySelectorAll('.tab-content').forEach(c => c.classList.remove('active'));
        tab.classList.add('active');
        document.getElementById(tab.getAttribute('data-tab')).classList.add('active');
    });
});

// -------------------- API BASE --------------------
const API_URL = "http://localhost:8080";

// -------------------- Helpers --------------------
async function fetchJSON(url, options = {}) {
    const res = await fetch(url, options);
    if (!res.ok) throw new Error(`Error ${res.status}`);
    return await res.json();
}

function renderTable(tableId, data, type) {
    const tbody = document.querySelector(`#${tableId} tbody`);
    tbody.innerHTML = "";
    data.forEach(item => {
        let row = document.createElement('tr');
        if (type === "dolls") {
            row.innerHTML = `
                <td>${item.nombre}</td>
                <td>${item.edad}</td>
                <td><span class="status ${item.activo ? 'status-sent' : 'status-draft'}">${item.activo ? 'Activa' : 'Inactiva'}</span></td>
                <td>${item.cartas}</td>
                <td class="action-buttons">
                    <button class="action-btn btn-primary" onclick="openEditDoll(${item.id})"><i class="fas fa-edit"></i></button>
                    <button class="action-btn btn-danger" onclick="deleteDoll(${item.id})"><i class="fas fa-trash"></i></button>
                </td>`;
        } else if (type === "clientes") {
            row.innerHTML = `
                <td>${item.nombre}</td>
                <td>${item.ciudad}</td>
                <td>${item.motivo}</td>
                <td>${item.contacto}</td>
                <td class="action-buttons">
                    <button class="action-btn btn-primary" onclick="openEditCliente(${item.id})"><i class="fas fa-edit"></i></button>
                    <button class="action-btn btn-danger" onclick="deleteCliente(${item.id})"><i class="fas fa-trash"></i></button>
                </td>`;
        } else if (type === "cartas") {
            row.innerHTML = `
                <td>${item.cliente_id}</td>
                <td>${item.doll_id}</td>
                <td>${item.fecha}</td>
                <td><span class="status ${item.estado}">${item.estado}</span></td>
                <td>${item.contenido}</td>
                <td class="action-buttons">
                    <button class="action-btn btn-primary" onclick="advanceCarta(${item.id}, '${item.estado}')"><i class="fas fa-forward"></i></button>
                    <button class="action-btn btn-danger" onclick="deleteCarta(${item.id})"><i class="fas fa-trash"></i></button>
                </td>`;
        }
        tbody.appendChild(row);
    });
    updateStats();
}

// -------------------- Load Data --------------------
async function loadDolls() {
    const dolls = await fetchJSON(`${API_URL}/dolls`);
    renderTable("dolls", dolls, "dolls");
}

async function loadClientes() {
    const clientes = await fetchJSON(`${API_URL}/clientes`);
    renderTable("clients", clientes, "clientes");
    // llenar select de clientes en formulario de cartas
    const select = document.getElementById('letterClient');
    select.innerHTML = "";
    clientes.forEach(c => {
        let option = document.createElement('option');
        option.value = c.id;
        option.textContent = `${c.nombre} (${c.ciudad})`;
        select.appendChild(option);
    });
}

async function loadCartas() {
    const cartas = await fetchJSON(`${API_URL}/cartas`);
    renderTable("letters", cartas, "cartas");
}

// -------------------- CRUD Dolls --------------------
async function createDoll(doll) {
    await fetchJSON(`${API_URL}/dolls`, {method: "POST", headers: {"Content-Type": "application/json"}, body: JSON.stringify(doll)});
    loadDolls();
}

async function updateDoll(id, doll) {
    await fetchJSON(`${API_URL}/dolls/${id}`, {method: "PUT", headers: {"Content-Type": "application/json"}, body: JSON.stringify(doll)});
    loadDolls();
}

async function deleteDoll(id) {
    if (confirm("¿Eliminar Doll?")) {
        await fetchJSON(`${API_URL}/dolls/${id}`, {method: "DELETE"});
        loadDolls();
    }
}

// -------------------- CRUD Clientes --------------------
async function createCliente(cliente) {
    await fetchJSON(`${API_URL}/clientes`, {method: "POST", headers: {"Content-Type": "application/json"}, body: JSON.stringify(cliente)});
    loadClientes();
}

async function updateCliente(id, cliente) {
    await fetchJSON(`${API_URL}/clientes/${id}`, {method: "PUT", headers: {"Content-Type": "application/json"}, body: JSON.stringify(cliente)});
    loadClientes();
}

async function deleteCliente(id) {
    if (confirm("¿Eliminar Cliente?")) {
        await fetchJSON(`${API_URL}/clientes/${id}`, {method: "DELETE"});
        loadClientes();
    }
}

// -------------------- CRUD Cartas --------------------
async function createCarta(carta) {
    await fetchJSON(`${API_URL}/cartas`, {method: "POST", headers: {"Content-Type": "application/json"}, body: JSON.stringify(carta)});
    loadCartas();
}

async function advanceCarta(id, estadoActual) {
    let nuevoEstado = null;
    if (estadoActual === "borrador") nuevoEstado = "revisado";
    else if (estadoActual === "revisado") nuevoEstado = "enviado";
    if (!nuevoEstado) return alert("No se puede avanzar más");
    await fetchJSON(`${API_URL}/cartas/${id}`, {method: "PUT", headers: {"Content-Type": "application/json"}, body: JSON.stringify({estado: nuevoEstado})});
    loadCartas();
}

async function deleteCarta(id) {
    if (confirm("¿Eliminar Carta?")) {
        try {
            await fetchJSON(`${API_URL}/cartas/${id}`, {method: "DELETE"});
            loadCartas();
        } catch (err) {
            alert("Solo se pueden eliminar cartas en borrador");
        }
    }
}

// -------------------- Formularios --------------------
const modals = {
    doll: document.getElementById('dollModal'),
    client: document.getElementById('clientModal'),
    letter: document.getElementById('letterModal')
};

function closeModal(type){ modals[type].style.display='none'; }

// Crear Doll
document.getElementById('dollForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const doll = {
        nombre: document.getElementById('dollName').value.trim(),
        edad: parseInt(document.getElementById('dollAge').value),
        activo: document.getElementById('dollStatus').value === 'active',
        cartas: 0
    };
    if (!doll.nombre || !doll.edad) return alert("Datos incompletos");
    await createDoll(doll);
    closeModal('doll');
});

// Crear Cliente
document.getElementById('clientForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const cliente = {
        nombre: document.getElementById('clientName').value.trim(),
        ciudad: document.getElementById('clientCity').value.trim(),
        motivo: document.getElementById('clientReason').value.trim(),
        contacto: document.getElementById('clientContact').value.trim()
    };
    if (!cliente.nombre || !cliente.ciudad || !cliente.motivo || !cliente.contacto) return alert("Datos incompletos");
    await createCliente(cliente);
    closeModal('client');
});

// Crear Carta
document.getElementById('letterForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const carta = {
        cliente_id: parseInt(document.getElementById('letterClient').value),
        fecha: document.getElementById('letterDate').value,
        contenido: document.getElementById('letterContent').value.trim()
    };
    if (!carta.cliente_id || !carta.fecha || !carta.contenido) return alert("Datos incompletos");
    await createCarta(carta);
    closeModal('letter');
});

// -------------------- Stats --------------------
function updateStats() {
    const dollCount = document.querySelectorAll('#dolls tbody tr').length;
    const clientCount = document.querySelectorAll('#clients tbody tr').length;
    const letterCount = document.querySelectorAll('#letters tbody tr').length;
    document.querySelector('.stat-card:nth-child(1) h3').textContent = dollCount;
    document.querySelector('.stat-card:nth-child(2) h3').textContent = clientCount;
    document.querySelector('.stat-card:nth-child(3) h3').textContent = letterCount;
}

// -------------------- Init --------------------
document.addEventListener('DOMContentLoaded', () => {
    loadDolls();
    loadClientes();
    loadCartas();
    document.getElementById('letterDate').valueAsDate = new Date();
});
