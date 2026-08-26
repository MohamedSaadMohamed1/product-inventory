package api

import (
	"net/http"
)

// ServeDashboard returns the HTML content for the interactive developer dashboard.
func ServeDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(dashboardHTML))
}

const dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Product Inventory & Stock Reservation Dashboard</title>
    <link href="https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;600;800&display=swap" rel="stylesheet">
    <style>
        :root {
            /* Light Theme: Yellow, White, Blue */
            --bg-light: #f4f6f9;        /* Soft gray-blue background */
            --bg-card: #ffffff;         /* Pure white for cards */
            --primary: #2563eb;         /* Classic royal blue */
            --primary-hover: #1d4ed8;   /* Darker blue */
            --primary-glow: rgba(37, 99, 235, 0.1);
            --accent-yellow: #facc15;   /* Bright yellow for highlights and primary branding */
            --accent-yellow-dark: #ca8a04;
            --success: #10b981;         /* Green */
            --danger: #ef4444;          /* Red */
            --warning: #f59e0b;         /* Orange-Yellow */
            --text-main: #1e293b;       /* Dark slate gray text */
            --text-muted: #64748b;      /* Muted slate text */
            --border: #e2e8f0;          /* Light border color */
            --bg-input: #f8fafc;        /* Input background */
        }

        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }

        body {
            font-family: 'Outfit', sans-serif;
            background-color: var(--bg-light);
            color: var(--text-main);
            line-height: 1.6;
            padding: 2rem;
            min-height: 100vh;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
        }

        header {
            text-align: center;
            margin-bottom: 3rem;
            position: relative;
        }

        h1 {
            font-size: 2.8rem;
            font-weight: 800;
            /* Blue and Yellow Gradient */
            background: linear-gradient(135deg, var(--primary) 0%, #3b82f6 50%, var(--accent-yellow-dark) 100%);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            margin-bottom: 0.5rem;
            letter-spacing: -0.03em;
        }

        .subtitle {
            color: var(--text-muted);
            font-size: 1.1rem;
            font-weight: 400;
        }

        /* Dashboard Grid Layout */
        .grid {
            display: grid;
            grid-template-columns: 2fr 1fr;
            gap: 2rem;
        }

        @media (max-width: 900px) {
            .grid {
                grid-template-columns: 1fr;
            }
        }

        .card {
            background-color: var(--bg-card);
            border: 1px solid var(--border);
            border-radius: 16px;
            padding: 1.8rem;
            box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
            transition: transform 0.2s, border-color 0.2s, box-shadow 0.2s;
        }

        .card:hover {
            border-color: var(--primary);
            box-shadow: 0 6px 25px rgba(37, 99, 235, 0.08);
        }

        .card-title {
            font-size: 1.3rem;
            font-weight: 600;
            margin-bottom: 1.5rem;
            display: flex;
            justify-content: space-between;
            align-items: center;
            border-bottom: 1px solid var(--border);
            padding-bottom: 0.8rem;
            color: var(--primary);
        }

        /* Products Grid */
        .products-container {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 1rem;
        }

        @media (max-width: 600px) {
            .products-container {
                grid-template-columns: 1fr;
            }
        }

        .product-card {
            background: var(--bg-light);
            border: 1px solid var(--border);
            border-radius: 12px;
            padding: 1.2rem;
            display: flex;
            flex-direction: column;
            justify-content: space-between;
            position: relative;
            overflow: hidden;
            transition: all 0.3s;
        }

        .product-card:hover {
            background: #ffffff;
            transform: translateY(-2px);
            border-color: var(--primary);
            box-shadow: 0 4px 15px rgba(0, 0, 0, 0.05);
        }

        .prod-id {
            position: absolute;
            top: 0.8rem;
            right: 0.8rem;
            font-size: 0.8rem;
            color: var(--text-muted);
            font-weight: 600;
        }

        .prod-name {
            font-size: 1.1rem;
            font-weight: 600;
            margin-bottom: 0.4rem;
            padding-right: 1.5rem;
            color: var(--text-main);
        }

        .prod-desc {
            font-size: 0.85rem;
            color: var(--text-muted);
            margin-bottom: 1rem;
            flex-grow: 1;
        }

        .prod-meta {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-top: 1rem;
            flex-wrap: wrap;
            gap: 0.5rem;
        }

        .prod-price {
            font-size: 1.2rem;
            font-weight: 800;
            color: var(--primary);
        }

        .stock-badge {
            padding: 0.3rem 0.6rem;
            border-radius: 20px;
            font-size: 0.75rem;
            font-weight: 600;
        }

        .stock-badge.in-stock {
            background: rgba(16, 185, 129, 0.1);
            color: var(--success);
            border: 1px solid rgba(16, 185, 129, 0.2);
        }

        .stock-badge.low-stock {
            /* Highlight low stock in Yellow theme style */
            background: rgba(250, 204, 21, 0.15);
            color: var(--accent-yellow-dark);
            border: 1px solid rgba(250, 204, 21, 0.3);
        }

        .stock-badge.out-stock {
            background: rgba(239, 68, 68, 0.1);
            color: var(--danger);
            border: 1px solid rgba(239, 68, 68, 0.2);
        }

        /* Order Panel & Logs */
        .form-group {
            margin-bottom: 1.2rem;
        }

        label {
            display: block;
            font-size: 0.85rem;
            color: var(--text-muted);
            margin-bottom: 0.5rem;
            font-weight: 600;
        }

        select, input {
            width: 100%;
            padding: 0.8rem;
            background-color: var(--bg-input);
            border: 1px solid var(--border);
            border-radius: 8px;
            color: var(--text-main);
            font-family: inherit;
            font-size: 0.95rem;
            outline: none;
            transition: border-color 0.2s, background-color 0.2s;
        }

        select:focus, input:focus {
            border-color: var(--primary);
            background-color: #ffffff;
        }

        .btn {
            display: block;
            width: 100%;
            padding: 0.9rem;
            /* Blue Gradient styled buttons */
            background: linear-gradient(135deg, var(--primary) 0%, #3b82f6 100%);
            border: none;
            border-radius: 8px;
            color: white;
            font-weight: 600;
            font-size: 1rem;
            cursor: pointer;
            transition: filter 0.2s, transform 0.1s;
            box-shadow: 0 4px 15px rgba(37, 99, 235, 0.2);
        }

        .btn:hover {
            filter: brightness(1.1);
        }

        .btn:active {
            transform: scale(0.98);
        }

        .console-container {
            margin-top: 1.5rem;
        }

        .console {
            /* Keep console dark for programming output contrast */
            background-color: #0f172a;
            border: 1px solid var(--border);
            border-radius: 8px;
            padding: 1rem;
            font-family: monospace;
            font-size: 0.85rem;
            max-height: 200px;
            overflow-y: auto;
            color: #38bdf8; /* Light blue output text */
        }

        .console .success {
            color: #4ade80; /* Light green */
        }

        .console .error {
            color: #f87171; /* Light red */
        }

        .console .info {
            color: #fef08a; /* Light yellow */
        }

        /* Order Log Row */
        .order-history-list {
            margin-top: 1.5rem;
            display: flex;
            flex-direction: column;
            gap: 0.8rem;
        }

        .order-item-log {
            background: var(--bg-light);
            border: 1px solid var(--border);
            border-radius: 8px;
            padding: 1rem;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }

        .order-item-log.cancelled {
            border-color: rgba(239, 68, 68, 0.2);
            background: rgba(239, 68, 68, 0.02);
        }

        .btn-cancel {
            background: transparent;
            border: 1px solid var(--danger);
            color: var(--danger);
            padding: 0.4rem 0.8rem;
            border-radius: 6px;
            font-size: 0.8rem;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.2s;
        }

        .btn-cancel:hover {
            background: var(--danger);
            color: white;
        }

        .cancelled-text {
            color: var(--danger);
            font-size: 0.8rem;
            font-weight: 600;
        }

        .btn-refresh {
            background: transparent;
            border: 1px solid var(--border);
            color: var(--text-muted);
            padding: 0.3rem 0.6rem;
            border-radius: 6px;
            cursor: pointer;
            font-size: 0.8rem;
            transition: all 0.2s;
        }

        .btn-refresh:hover {
            border-color: var(--primary);
            color: var(--primary);
        }

        /* Auth Panel styling */
        #authLoggedIn {
            background: rgba(250, 204, 21, 0.05);
            border: 1px dashed var(--accent-yellow-dark);
            padding: 1rem;
            border-radius: 8px;
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>Product Inventory Dashboard</h1>
            <div class="subtitle">Interactive stock reservation and order lifecycle tester (JWT Protected)</div>
        </header>

        <div class="grid">
            <!-- Left Side: Products Catalog -->
            <div class="card">
                <div class="card-title">
                    <span>Products Catalog</span>
                    <button class="btn-refresh" onclick="fetchProducts()">Refresh Catalog</button>
                </div>
                <div class="products-container" id="productsGrid">
                    <!-- Products will be injected here -->
                </div>
            </div>

            <!-- Right Side: Auth, Order Panel & Logs -->
            <div>
                <!-- Authentication Card -->
                <div class="card" style="margin-bottom: 1.5rem;" id="authCard">
                    <div class="card-title">Authentication</div>
                    
                    <div id="authLoggedOut">
                        <div class="form-group">
                            <label for="authEmail">Email</label>
                            <input type="email" id="authEmail" placeholder="e.g. customer@example.com or admin@example.com" required>
                        </div>
                        <div class="form-group">
                            <label for="authPassword">Password</label>
                            <input type="password" id="authPassword" placeholder="e.g. customerpassword or adminpassword" required>
                        </div>
                        <div class="form-group" style="display: none;" id="roleGroup">
                            <label for="authRole">Role</label>
                            <select id="authRole">
                                <option value="customer">Customer</option>
                                <option value="admin">Admin</option>
                            </select>
                        </div>
                        <div style="display: flex; gap: 0.5rem; margin-top: 1rem;">
                            <button class="btn" style="flex: 1;" onclick="login(event)" id="btnLogin">Login</button>
                            <button class="btn" style="flex: 1; background: var(--text-muted);" onclick="toggleAuthMode(event)" id="btnRegToggle">Register</button>
                        </div>
                    </div>

                    <div id="authLoggedIn" style="display: none;">
                        <div style="display: flex; justify-content: space-between; align-items: center;">
                            <div>
                                <div style="font-weight: 800; font-size: 1.1rem; color: var(--primary);" id="loggedInUser">Guest</div>
                                <div style="font-size: 0.8rem; color: var(--text-muted); text-transform: uppercase;" id="loggedInRole">Role: customer</div>
                            </div>
                            <button class="btn-cancel" onclick="logout(event)">Logout</button>
                        </div>
                    </div>
                </div>

                <!-- Admin panel (Hidden by default, displayed only when user is admin) -->
                <div class="card" id="adminPanel" style="display: none; margin-bottom: 1.5rem;">
                    <div class="card-title">Admin: Manage Products</div>
                    
                    <div style="display: flex; border-bottom: 1px solid var(--border); margin-bottom: 1.2rem; gap: 0.5rem;">
                        <button id="tabAddBtn" class="btn-refresh" style="border: none; border-bottom: 2px solid var(--primary); border-radius: 0; padding: 0.5rem 1rem; color: var(--primary); font-weight: 600;" onclick="switchAdminTab('add')">Add Product</button>
                        <button id="tabUpdateBtn" class="btn-refresh" style="border: none; border-bottom: 2px solid transparent; border-radius: 0; padding: 0.5rem 1rem; color: var(--text-muted);" onclick="switchAdminTab('update')">Update Product</button>
                    </div>

                    <!-- Add Product Form -->
                    <form id="addProductForm" onsubmit="addProduct(event)">
                        <div class="form-group">
                            <label for="newProdName">Product Name</label>
                            <input type="text" id="newProdName" placeholder="e.g. Sony PlayStation 5 Pro" required>
                        </div>
                        <div class="form-group">
                            <label for="newProdDesc">Description</label>
                            <input type="text" id="newProdDesc" placeholder="e.g. 1TB SSD Ultra HD Console" required>
                        </div>
                        <div style="display: flex; gap: 0.5rem;">
                            <div class="form-group" style="flex: 1;">
                                <label for="newProdPrice">Price (USD)</label>
                                <input type="number" id="newProdPrice" step="0.01" min="0.01" placeholder="499.99" required>
                            </div>
                            <div class="form-group" style="flex: 1;">
                                <label for="newProdStock">Stock</label>
                                <input type="number" id="newProdStock" min="0" value="10" required>
                            </div>
                        </div>
                        <button type="submit" class="btn">Create Product</button>
                    </form>

                    <!-- Update Product Form -->
                    <form id="updateProductForm" onsubmit="updateProduct(event)" style="display: none;">
                        <div class="form-group">
                            <label for="updateProdSelect">Select Product to Edit</label>
                            <select id="updateProdSelect" onchange="prefillUpdateForm()" required>
                                <!-- Options injected dynamically -->
                            </select>
                        </div>
                        <div class="form-group">
                            <label for="updateProdName">Product Name</label>
                            <input type="text" id="updateProdName" required>
                        </div>
                        <div class="form-group">
                            <label for="updateProdDesc">Description</label>
                            <input type="text" id="updateProdDesc" required>
                        </div>
                        <div style="display: flex; gap: 0.5rem;">
                            <div class="form-group" style="flex: 1;">
                                <label for="updateProdPrice">Price (USD)</label>
                                <input type="number" id="updateProdPrice" step="0.01" min="0.01" required>
                            </div>
                            <div class="form-group" style="flex: 1;">
                                <label for="updateProdStock">Stock Quantity</label>
                                <input type="number" id="updateProdStock" min="0" required>
                            </div>
                        </div>
                        <button type="submit" class="btn" style="background: linear-gradient(135deg, var(--warning) 0%, #d97706 100%); box-shadow: 0 4px 15px rgba(245, 158, 11, 0.2);">Save Changes</button>
                    </form>
                </div>

                <!-- Ordering Card (Now supports multi-product rows) -->
                <div class="card">
                    <div class="card-title">Reserve Stock</div>
                    <form id="orderForm" onsubmit="createOrder(event)">
                        <div id="orderItemsContainer">
                            <!-- Dynamic product rows go here -->
                        </div>
                        
                        <button type="button" class="btn-refresh" style="width: 100%; margin-bottom: 1.5rem; border-color: var(--primary); color: var(--primary); font-weight: 600;" onclick="addProductRow()">+ Add Another Product</button>
                        
                        <button type="submit" class="btn">Place Order</button>
                    </form>

                    <div class="console-container">
                        <label>Execution Output</label>
                        <div class="console" id="consoleOutput">Waiting for actions...</div>
                    </div>
                </div>

                <!-- Orders Log Card -->
                <div class="card" style="margin-top: 1.5rem;">
                    <div class="card-title">Orders Log</div>
                    <div class="order-history-list" id="ordersList">
                        <!-- Placed orders will be listed here -->
                    </div>
                </div>
            </div>
        </div>
    </div>

    <script>
        let products = [];
        let orders = [];
        let isRegisterMode = false;

        // Auth states
        let token = localStorage.getItem('jwt_token') || '';
        let userRole = localStorage.getItem('user_role') || '';
        let email = localStorage.getItem('email') || '';

        function checkAuthStatus() {
            const loggedOutDiv = document.getElementById('authLoggedOut');
            const loggedInDiv = document.getElementById('authLoggedIn');
            const adminPanelDiv = document.getElementById('adminPanel');
            
            if (token) {
                loggedOutDiv.style.display = 'none';
                loggedInDiv.style.display = 'block';
                document.getElementById('loggedInUser').innerText = email;
                document.getElementById('loggedInRole').innerText = 'Role: ' + userRole;

                if (userRole === 'admin') {
                    adminPanelDiv.style.display = 'block';
                } else {
                    adminPanelDiv.style.display = 'none';
                }
            } else {
                loggedOutDiv.style.display = 'block';
                loggedInDiv.style.display = 'none';
                adminPanelDiv.style.display = 'none';
            }
            fetchOrders();
        }

        function toggleAuthMode(event) {
            event.preventDefault();
            isRegisterMode = !isRegisterMode;
            const roleGroup = document.getElementById('roleGroup');
            const btnLogin = document.getElementById('btnLogin');
            const btnRegToggle = document.getElementById('btnRegToggle');

            if (isRegisterMode) {
                roleGroup.style.display = 'block';
                btnLogin.innerText = 'Sign Up';
                btnRegToggle.innerText = 'Go to Login';
            } else {
                roleGroup.style.display = 'none';
                btnLogin.innerText = 'Login';
                btnRegToggle.innerText = 'Register';
            }
        }

        async function login(event) {
            event.preventDefault();
            const emailVal = document.getElementById('authEmail').value;
            const passVal = document.getElementById('authPassword').value;
            const roleVal = document.getElementById('authRole').value;

            if (!emailVal || !passVal) {
                logConsole('Please enter email and password', 'error');
                return;
            }

            const path = isRegisterMode ? '/api/v1/auth/register' : '/api/v1/auth/login';
            const bodyObj = isRegisterMode 
                ? { email: emailVal, password: passVal, role: roleVal } 
                : { email: emailVal, password: passVal };

            logConsole((isRegisterMode ? 'Registering...' : 'Logging in...') + ' Email: ' + emailVal, 'info');

            try {
                const response = await fetch(path, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(bodyObj)
                });

                const data = await response.json();

                if (response.ok) {
                    token = data.token;
                    userRole = data.role;
                    email = data.email;

                    localStorage.setItem('jwt_token', token);
                    localStorage.setItem('user_role', userRole);
                    localStorage.setItem('email', email);

                    logConsole('Successfully authenticated as ' + email + ' (' + userRole + ')', 'success');
                    
                    // Clear inputs
                    document.getElementById('authEmail').value = '';
                    document.getElementById('authPassword').value = '';
                    
                    checkAuthStatus();
                    await fetchProducts();
                } else {
                    const errMsg = data.error ? data.error.message : 'Auth failed';
                    logConsole('Auth error: ' + errMsg, 'error');
                }
            } catch (err) {
                logConsole('Auth connection error: ' + err.message, 'error');
            }
        }

        function logout(event) {
            event.preventDefault();
            token = '';
            userRole = '';
            email = '';

            localStorage.removeItem('jwt_token');
            localStorage.removeItem('user_role');
            localStorage.removeItem('email');

            logConsole('Logged out successfully', 'info');
            checkAuthStatus();
            fetchProducts();
        }

        async function fetchProducts() {
            try {
                const headersObj = {};
                if (token) {
                    headersObj['Authorization'] = 'Bearer ' + token;
                }

                const response = await fetch('/api/v1/products', { headers: headersObj });
                if (!response.ok) throw new Error('Failed to load products');
                products = await response.json();
                renderProducts();
                refreshAllProductSelects();
                populateUpdateDropdown();
            } catch (err) {
                logConsole('Error fetching products: ' + err.message, 'error');
            }
        }

        function renderProducts() {
            const grid = document.getElementById('productsGrid');
            grid.innerHTML = '';

            products.forEach(p => {
                let badgeClass = 'in-stock';
                let badgeText = 'In Stock';

                if (p.stock_quantity === 0) {
                    badgeClass = 'out-stock';
                    badgeText = 'Out of Stock';
                } else if (p.stock_quantity < 5) {
                    badgeClass = 'low-stock';
                    badgeText = 'Low Stock (' + p.stock_quantity + ')';
                } else {
                    badgeText = 'In Stock (' + p.stock_quantity + ')';
                }

                const priceFormatted = (p.price / 100).toLocaleString('en-US', {
                    style: 'currency',
                    currency: 'USD'
                });

                let deleteBtnHtml = '';
                if (userRole === 'admin') {
                    deleteBtnHtml = '<button class="btn-cancel" style="border-color: var(--danger); color: var(--danger); padding: 0.2rem 0.5rem; font-size: 0.75rem;" onclick="deleteProduct(' + p.id + ')">Delete</button>';
                }

                grid.innerHTML += 
                    '<div class="product-card">' +
                        '<span class="prod-id">#' + p.id + '</span>' +
                        '<div class="prod-name">' + escapeHTML(p.name) + '</div>' +
                        '<div class="prod-desc">' + escapeHTML(p.description) + '</div>' +
                        '<div class="prod-meta">' +
                            '<span class="prod-price">' + priceFormatted + '</span>' +
                            '<span class="stock-badge ' + badgeClass + '">' + badgeText + '</span>' +
                            deleteBtnHtml +
                        '</div>' +
                    '</div>';
            });
        }

        // Multi-Product Form Logic
        function addProductRow() {
            const container = document.getElementById('orderItemsContainer');
            const rowId = 'row_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);

            const rowDiv = document.createElement('div');
            rowDiv.className = 'order-item-row';
            rowDiv.id = rowId;
            rowDiv.style.display = 'flex';
            rowDiv.style.gap = '0.5rem';
            rowDiv.style.alignItems = 'flex-end';
            rowDiv.style.marginBottom = '0.8rem';

            rowDiv.innerHTML = 
                '<div style="flex: 2;">' +
                    '<label>Product</label>' +
                    '<select class="prod-select" required></select>' +
                '</div>' +
                '<div style="flex: 1; min-width: 80px;">' +
                    '<label>Qty</label>' +
                    '<input type="number" class="qty-input" min="1" value="1" required>' +
                '</div>' +
                '<div>' +
                    '<button type="button" class="btn-cancel" style="padding: 0.75rem; border-color: var(--danger); color: var(--danger);" onclick="removeProductRow(\'' + rowId + '\')">Remove</button>' +
                '</div>';

            container.appendChild(rowDiv);
            populateProductSelect(rowDiv.querySelector('.prod-select'));
        }

        function removeProductRow(rowId) {
            const row = document.getElementById(rowId);
            if (row) {
                row.remove();
            }
            // Ensure there is always at least one row
            const container = document.getElementById('orderItemsContainer');
            if (container.children.length === 0) {
                addProductRow();
            }
        }

        function populateProductSelect(selectElement) {
            selectElement.innerHTML = '';
            products.forEach(p => {
                selectElement.innerHTML += 
                    '<option value="' + p.id + '" ' + (p.stock_quantity === 0 ? 'disabled' : '') + '>' +
                        escapeHTML(p.name) + ' (' + p.stock_quantity + ' available)' +
                    '</option>';
            });
        }

        function refreshAllProductSelects() {
            const selects = document.querySelectorAll('.prod-select');
            selects.forEach(select => {
                const currentVal = select.value;
                populateProductSelect(select);
                if (currentVal) {
                    select.value = currentVal;
                }
            });
        }

        // Admin Management functions
        function switchAdminTab(tab) {
            const formAdd = document.getElementById('addProductForm');
            const formUpdate = document.getElementById('updateProductForm');
            const tabAddBtn = document.getElementById('tabAddBtn');
            const tabUpdateBtn = document.getElementById('tabUpdateBtn');

            if (tab === 'add') {
                formAdd.style.display = 'block';
                formUpdate.style.display = 'none';
                tabAddBtn.style.borderBottom = '2px solid var(--primary)';
                tabAddBtn.style.color = 'var(--primary)';
                tabUpdateBtn.style.borderBottom = '2px solid transparent';
                tabUpdateBtn.style.color = 'var(--text-muted)';
            } else {
                formAdd.style.display = 'none';
                formUpdate.style.display = 'block';
                tabAddBtn.style.borderBottom = '2px solid transparent';
                tabAddBtn.style.color = 'var(--text-muted)';
                tabUpdateBtn.style.borderBottom = '2px solid var(--primary)';
                tabUpdateBtn.style.color = 'var(--primary)';
                prefillUpdateForm();
            }
        }

        function populateUpdateDropdown() {
            const select = document.getElementById('updateProdSelect');
            if (!select) return;
            select.innerHTML = '';
            products.forEach(p => {
                select.innerHTML += '<option value="' + p.id + '">' + escapeHTML(p.name) + '</option>';
            });
        }

        function prefillUpdateForm() {
            const selectVal = parseInt(document.getElementById('updateProdSelect').value);
            if (!selectVal) return;
            const p = products.find(prod => prod.id === selectVal);
            if (!p) return;

            document.getElementById('updateProdName').value = p.name;
            document.getElementById('updateProdDesc').value = p.description;
            document.getElementById('updateProdPrice').value = (p.price / 100).toFixed(2);
            document.getElementById('updateProdStock').value = p.stock_quantity;
        }

        async function addProduct(event) {
            event.preventDefault();
            if (!token || userRole !== 'admin') {
                logConsole('Error: Admin permissions required.', 'error');
                return;
            }

            const name = document.getElementById('newProdName').value;
            const desc = document.getElementById('newProdDesc').value;
            const price = parseFloat(document.getElementById('newProdPrice').value);
            const stock = parseInt(document.getElementById('newProdStock').value);

            const priceCents = Math.round(price * 100);

            logConsole('Creating product: ' + name + '...', 'info');

            try {
                const response = await fetch('/api/v1/products', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': 'Bearer ' + token
                    },
                    body: JSON.stringify({
                        name: name,
                        description: desc,
                        price: priceCents,
                        stock_quantity: stock
                    })
                });

                const data = await response.json();

                if (response.ok) {
                    logConsole('Success: Created product #' + data.id + ' (' + data.name + ')', 'success');
                    // Reset inputs
                    document.getElementById('newProdName').value = '';
                    document.getElementById('newProdDesc').value = '';
                    document.getElementById('newProdPrice').value = '';
                    document.getElementById('newProdStock').value = '10';
                    await fetchProducts();
                } else {
                    const errMsg = data.error ? data.error.message : 'Failed to create product';
                    logConsole('Create Product failed: ' + errMsg, 'error');
                }
            } catch (err) {
                logConsole('Connection error: ' + err.message, 'error');
            }
        }

        async function updateProduct(event) {
            event.preventDefault();
            if (!token || userRole !== 'admin') {
                logConsole('Error: Admin permissions required.', 'error');
                return;
            }

            const productId = parseInt(document.getElementById('updateProdSelect').value);
            const name = document.getElementById('updateProdName').value;
            const desc = document.getElementById('updateProdDesc').value;
            const price = parseFloat(document.getElementById('updateProdPrice').value);
            const stock = parseInt(document.getElementById('updateProdStock').value);

            const priceCents = Math.round(price * 100);

            logConsole('Updating product #' + productId + '...', 'info');

            try {
                const response = await fetch('/api/v1/products/' + productId, {
                    method: 'PUT',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': 'Bearer ' + token
                    },
                    body: JSON.stringify({
                        name: name,
                        description: desc,
                        price: priceCents,
                        stock_quantity: stock
                    })
                });

                const data = await response.json();

                if (response.ok) {
                    logConsole('Success: Updated product #' + data.id, 'success');
                    await fetchProducts();
                } else {
                    const errMsg = data.error ? data.error.message : 'Failed to update product';
                    logConsole('Update Product failed: ' + errMsg, 'error');
                }
            } catch (err) {
                logConsole('Connection error: ' + err.message, 'error');
            }
        }

        async function deleteProduct(productId) {
            if (!token || userRole !== 'admin') {
                logConsole('Error: Admin permissions required.', 'error');
                return;
            }

            if (!confirm('Are you sure you want to delete Product #' + productId + '?')) {
                return;
            }

            logConsole('Deleting product #' + productId + '...', 'info');

            try {
                const response = await fetch('/api/v1/products/' + productId, {
                    method: 'DELETE',
                    headers: {
                        'Authorization': 'Bearer ' + token
                    }
                });

                if (response.status === 204 || response.ok) {
                    logConsole('Success: Deleted product #' + productId, 'success');
                    await fetchProducts();
                } else {
                    const data = await response.json();
                    const errMsg = data.error ? data.error.message : 'Failed to delete product';
                    logConsole('Delete Product failed: ' + errMsg, 'error');
                }
            } catch (err) {
                logConsole('Connection error: ' + err.message, 'error');
            }
        }

        async function fetchOrders() {
            if (!token) {
                orders = [];
                renderOrders();
                return;
            }
            try {
                const response = await fetch('/api/v1/orders?page=1&limit=50', {
                    headers: { 'Authorization': 'Bearer ' + token }
                });
                if (response.ok) {
                    orders = await response.json();
                    renderOrders();
                } else {
                    logConsole('Failed to load orders history', 'error');
                }
            } catch (err) {
                logConsole('Error loading orders: ' + err.message, 'error');
            }
        }

        async function createOrder(event) {
            event.preventDefault();
            if (!token) {
                logConsole('Error: You must be logged in to place an order.', 'error');
                return;
            }

            // Gather all items from the dynamic rows
            const items = [];
            const rows = document.querySelectorAll('.order-item-row');
            rows.forEach(row => {
                const prodId = parseInt(row.querySelector('.prod-select').value);
                const qty = parseInt(row.querySelector('.qty-input').value);
                if (prodId && qty > 0) {
                    items.push({ product_id: prodId, quantity: qty });
                }
            });

            if (items.length === 0) {
                logConsole('Error: Please add at least one product row with valid quantity.', 'error');
                return;
            }

            logConsole('Sending order request for ' + items.length + ' items...', 'info');

            try {
                const response = await fetch('/api/v1/orders', {
                    method: 'POST',
                    headers: { 
                        'Content-Type': 'application/json',
                        'Authorization': 'Bearer ' + token
                    },
                    body: JSON.stringify({
                        items: items
                    })
                });

                const data = await response.json();

                if (response.ok) {
                    logConsole('Success: Order #' + data.id + ' created! Total: $' + (data.total_amount / 100).toFixed(2), 'success');
                    
                    // Reset order items form rows back to one default row
                    document.getElementById('orderItemsContainer').innerHTML = '';
                    addProductRow();
                    
                    orders.unshift(data); // Add to beginning
                    renderOrders();
                    await fetchProducts();
                } else {
                    const errMsg = data.error ? data.error.message : 'Unknown error';
                    const errCode = data.error ? data.error.code : 'ERROR';
                    logConsole('Failed (' + errCode + '): ' + errMsg, 'error');
                }
            } catch (err) {
                logConsole('Connection error: ' + err.message, 'error');
            }
        }

        async function cancelOrder(orderId) {
            if (!token) {
                logConsole('Error: You must be logged in to cancel an order.', 'error');
                return;
            }

            logConsole('Sending cancellation request for Order #' + orderId + '...', 'info');

            try {
                const response = await fetch('/api/v1/orders/' + orderId + '/cancel', {
                    method: 'POST',
                    headers: {
                        'Authorization': 'Bearer ' + token
                    }
                });

                const data = await response.json();

                if (response.ok) {
                    logConsole('Success: Order #' + orderId + ' cancelled! Stock restored.', 'success');
                    
                    // Update local state status
                    const o = orders.find(ord => ord.id === orderId);
                    if (o) {
                        o.status = 'CANCELLED';
                    }
                    renderOrders();
                    await fetchProducts();
                } else {
                    const errMsg = data.error ? data.error.message : 'Unknown error';
                    logConsole('Failed: ' + errMsg, 'error');
                }
            } catch (err) {
                logConsole('Connection error: ' + err.message, 'error');
            }
        }

        function renderOrders() {
            const list = document.getElementById('ordersList');
            list.innerHTML = '';

            if (orders.length === 0) {
                list.innerHTML = '<div style="text-align: center; color: var(--text-muted); font-size: 0.9rem;">No orders placed in this session</div>';
                return;
            }

            orders.forEach(o => {
                const priceFormatted = (o.total_amount / 100).toLocaleString('en-US', {
                    style: 'currency',
                    currency: 'USD'
                });

                const isCancelled = o.status === 'CANCELLED';

                // Display order item breakdown
                let itemsBreakdown = '';
                if (o.items && o.items.length > 0) {
                    itemsBreakdown += '<div style="font-size: 0.75rem; color: var(--text-muted); margin-top: 0.3rem;">Items: ';
                    const itemsDesc = o.items.map(it => 'Prod #' + it.product_id + ' (x' + it.quantity + ')');
                    itemsBreakdown += itemsDesc.join(', ') + '</div>';
                }

                list.innerHTML += 
                    '<div class="order-history-list">' +
                        '<div class="order-item-log ' + (isCancelled ? 'cancelled' : '') + '">' +
                            '<div>' +
                                '<div style="font-weight: 600; font-size: 0.95rem;">Order #' + o.id + '</div>' +
                                '<div style="font-size: 0.8rem; font-weight: 600; color: var(--primary);">Total: ' + priceFormatted + '</div>' +
                                itemsBreakdown +
                            '</div>' +
                            '<div>' +
                                (isCancelled 
                                    ? '<span class="cancelled-text">CANCELLED</span>' 
                                    : '<button class="btn-cancel" onclick="cancelOrder(' + o.id + ')">Cancel Order</button>'
                                ) +
                            '</div>' +
                        '</div>' +
                    '</div>';
            });
        }

        function logConsole(message, type = 'info') {
            const consoleBox = document.getElementById('consoleOutput');
            const timestamp = new Date().toLocaleTimeString();
            const logClass = type === 'error' ? 'error' : (type === 'success' ? 'success' : 'info');
            
            if (consoleBox.innerHTML === 'Waiting for actions...') {
                consoleBox.innerHTML = '';
            }

            consoleBox.innerHTML += '<div class="' + logClass + '">[' + timestamp + '] ' + escapeHTML(message) + '</div>';
            consoleBox.scrollTop = consoleBox.scrollHeight;

            // Also log to the browser developer console
            if (type === 'error') {
                console.error('[' + timestamp + '] ' + message);
            } else if (type === 'success') {
                console.log('%c[' + timestamp + '] ' + message, 'color: #10b981; font-weight: bold;');
            } else {
                console.log('[' + timestamp + '] ' + message);
            }
        }

        function escapeHTML(str) {
            return str.replace(/[&<>'"]/g, 
                tag => ({
                    '&': '&amp;',
                    '<': '&lt;',
                    '>': '&gt;',
                    "'": '&#39;',
                    '"': '&quot;'
                }[tag] || tag)
            );
        }

        // Initialize dashboard
        checkAuthStatus();
        fetchProducts();
        
        // Add first order item row by default
        addProductRow();
    </script>
</body>
</html>
`
