// script.js from Java sample
const apiBaseUrl = '/api/item'; // Assuming relative path works as server serves both

// DOM Elements
const searchBtn = document.getElementById('search-btn');
const searchInput = document.getElementById('search-input');
const retrieveRandomBtn = document.getElementById('retrieve-random-btn');
const createBtn = document.getElementById('create-btn');

const nameInput = document.getElementById('name');
const descInput = document.getElementById('description');
const priceInput = document.getElementById('price');

const errorMsg = document.getElementById('error-msg');
const successMsg = document.getElementById('success-msg');
const itemsContainer = document.getElementById('items');

const fromCacheSpan = document.getElementById('from-cache');
const timeToFetchSpan = document.getElementById('time-to-fetch');

// Helpers
function showMessage(type, msg) {
    if (type === 'error') {
        errorMsg.textContent = msg;
        errorMsg.classList.remove('hidden');
        successMsg.classList.add('hidden');
    } else {
        successMsg.textContent = msg;
        successMsg.classList.remove('hidden');
        errorMsg.classList.add('hidden');
    }
    setTimeout(() => {
        errorMsg.classList.add('hidden');
        successMsg.classList.add('hidden');
    }, 5000);
}

function renderItem(item) {
    return `
    <div class="p-4 bg-white rounded shadow border border-gray-200">
      <div class="flex justify-between items-start">
        <div>
          <h3 class="text-lg font-bold">${item.name} <span class="text-sm font-normal text-gray-500">#${item.id}</span></h3>
          <p class="text-gray-700 mt-1">${item.description}</p>
          <p class="text-green-600 font-semibold mt-2">$${item.price}</p>
        </div>
        <button onclick="deleteItem(${item.id})" class="text-red-500 hover:text-red-700">Delete</button>
      </div>
    </div>
  `;
}

function updateAnalysis(fromCache, time) {
    fromCacheSpan.textContent = fromCache;
    timeToFetchSpan.textContent = time + 'ms';
}

// Actions
async function searchItem() {
    const id = searchInput.value;
    if (!id) return;

    itemsContainer.innerHTML = '<p class="text-center text-gray-500">Loading...</p>';

    const start = Date.now();
    try {
        const res = await fetch(`${apiBaseUrl}/${id}`);
        if (!res.ok) throw new Error('Item not found');
        const item = await res.json();
        const end = Date.now();

        itemsContainer.innerHTML = renderItem(item);
        updateAnalysis(item.fromCache, end - start);
    } catch (err) {
        itemsContainer.innerHTML = '';
        showMessage('error', err.message);
        updateAnalysis('null', 'null');
    }
}

async function getRandom() {
    itemsContainer.innerHTML = '<p class="text-center text-gray-500">Loading...</p>';
    // Note: Java sample returns {items: [...]}. We should handle that.
    try {
        const res = await fetch(`${apiBaseUrl}/random`);
        if (!res.ok) throw new Error('Failed to fetch random items');
        const data = await res.json();

        if (data.items && data.items.length > 0) {
            itemsContainer.innerHTML = data.items.map(renderItem).join('');
        } else {
            itemsContainer.innerHTML = '<p class="text-center text-gray-500">No items found.</p>';
        }
        updateAnalysis('N/A', 'N/A');
    } catch (err) {
        showMessage('error', err.message);
    }
}

async function createItem() {
    const name = nameInput.value;
    const desc = descInput.value;
    const price = parseFloat(priceInput.value);

    if (!name || !desc || isNaN(price)) {
        showMessage('error', 'Please fill all fields correctly');
        return;
    }

    try {
        const res = await fetch(`${apiBaseUrl}/create`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name, description: desc, price })
        });

        if (!res.ok) throw new Error('Failed to create item');
        const data = await res.json();

        showMessage('success', `Item created with ID: ${data.id}`);
        nameInput.value = '';
        descInput.value = '';
        priceInput.value = '';
    } catch (err) {
        showMessage('error', err.message);
    }
}

async function deleteItem(id) {
    if (!confirm('Are you sure?')) return;
    try {
        const res = await fetch(`${apiBaseUrl}/delete/${id}`, { method: 'DELETE' });
        if (!res.ok) throw new Error('Failed to delete');
        showMessage('success', `Item ${id} deleted`);
        itemsContainer.innerHTML = '';
    } catch (err) {
        showMessage('error', err.message);
    }
}

// Events
searchBtn.addEventListener('click', searchItem);
retrieveRandomBtn.addEventListener('click', getRandom);
createBtn.addEventListener('click', createItem);
window.deleteItem = deleteItem; // Expose to global scope for inline onclick
