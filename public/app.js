const API_BASE = '/api';
const state = {
    skills: [],
    filteredSkills: [],
    bundle: [],
    filters: {
        category: '',
        sfia: '',
        framework: '',
        search: ''
    }
};

document.addEventListener('DOMContentLoaded', () => {
    initEventListeners();
    loadSkills();
});

function initEventListeners() {
    const searchInput = document.getElementById('searchInput');
    const categoryFilter = document.getElementById('categoryFilter');
    const sfiaFilter = document.getElementById('sfiaFilter');
    const frameworkFilter = document.getElementById('frameworkFilter');
    const clearFilters = document.getElementById('clearFilters');
    const closeModal = document.getElementById('closeModal');
    const closeBundle = document.getElementById('closeBundle');
    const bundleToggle = document.getElementById('bundleToggle');
    const downloadZip = document.getElementById('downloadZip');
    const modal = document.getElementById('skillModal');

    searchInput.addEventListener('input', debounce(handleSearch, 300));
    categoryFilter.addEventListener('change', handleFilterChange);
    sfiaFilter.addEventListener('change', handleFilterChange);
    frameworkFilter.addEventListener('change', handleFilterChange);
    clearFilters.addEventListener('click', clearAllFilters);

    closeModal.addEventListener('click', () => closeModalFn(modal));
    modal.querySelector('.modal-backdrop').addEventListener('click', () => closeModalFn(modal));

    closeBundle.addEventListener('click', () => toggleBundle(false));
    bundleToggle.addEventListener('click', () => toggleBundle());
    downloadZip.addEventListener('click', downloadBundleZip);

    document.addEventListener('keydown', (e) => {
        if (e.key === '/' && document.activeElement !== searchInput) {
            e.preventDefault();
            searchInput.focus();
        }
        if (e.key === 'Escape') {
            closeModalFn(modal);
            toggleBundle(false);
        }
    });
}

async function loadSkills() {
    try {
        const res = await fetch(`${API_BASE}/skills`);
        const data = await res.json();
        state.skills = data.data || [];
        await loadCategories();
        applyFilters();
    } catch (err) {
        showToast('Failed to load skills');
        console.error(err);
    }
}

async function loadCategories() {
    try {
        const res = await fetch(`${API_BASE}/categories`);
        const data = await res.json();
        const select = document.getElementById('categoryFilter');
        
        const currentValue = select.value;
        select.innerHTML = '<option value="">All Categories</option>';
        
        (data.data || []).forEach(cat => {
            const option = document.createElement('option');
            option.value = cat.name;
            option.textContent = `${cat.name} (${cat.count})`;
            select.appendChild(option);
        });
        
        select.value = currentValue;
    } catch (err) {
        console.error('Failed to load categories:', err);
    }
}

function handleSearch(e) {
    state.filters.search = e.target.value;
    applyFilters();
}

function handleFilterChange(e) {
    const filterType = e.target.id.replace('Filter', '');
    state.filters[filterType] = e.target.value;
    applyFilters();
}

function clearAllFilters() {
    state.filters = {
        category: '',
        sfia: '',
        framework: '',
        search: ''
    };
    document.getElementById('searchInput').value = '';
    document.getElementById('categoryFilter').value = '';
    document.getElementById('sfiaFilter').value = '';
    document.getElementById('frameworkFilter').value = '';
    applyFilters();
}

function applyFilters() {
    let filtered = [...state.skills];

    if (state.filters.search) {
        const query = state.filters.search.toLowerCase();
        filtered = filtered.filter(skill =>
            skill.name.toLowerCase().includes(query) ||
            skill.description.toLowerCase().includes(query) ||
            skill.tags.some(tag => tag.toLowerCase().includes(query))
        );
    }

    if (state.filters.category) {
        filtered = filtered.filter(skill => 
            skill.category.toLowerCase() === state.filters.category.toLowerCase()
        );
    }

    if (state.filters.sfia) {
        const level = parseInt(state.filters.sfia);
        filtered = filtered.filter(skill => 
            skill.sfia && skill.sfia.level === level
        );
    }

    if (state.filters.framework) {
        filtered = filtered.filter(skill => 
            skill.framework && skill.framework.type.toLowerCase() === state.filters.framework.toLowerCase()
        );
    }

    state.filteredSkills = filtered;
    renderSkills();
}

function renderSkills() {
    const grid = document.getElementById('skillsGrid');
    const emptyState = document.getElementById('emptyState');
    const stats = document.getElementById('skillCount');

    stats.textContent = state.filteredSkills.length;

    if (state.filteredSkills.length === 0) {
        grid.innerHTML = '';
        emptyState.style.display = 'block';
        return;
    }

    emptyState.style.display = 'none';
    grid.innerHTML = state.filteredSkills.map(skill => createSkillCard(skill)).join('');

    grid.querySelectorAll('.skill-card').forEach(card => {
        card.addEventListener('click', () => openSkillModal(card.dataset.slug));
    });
}

function createSkillCard(skill) {
    const sfiaBadge = skill.sfia ? 
        `<span class="badge sfia">SFIA L${skill.sfia.level}</span>` : '';
    const frameworkBadge = skill.framework ? 
        `<span class="badge framework">${skill.framework.type.toUpperCase()}</span>` : '';
    const categoryBadge = skill.category ? 
        `<span class="badge category">${skill.category}</span>` : '';

    return `
        <div class="skill-card" data-slug="${skill.slug}">
            <div class="skill-card-header">
                <h3 class="skill-card-title">${escapeHtml(skill.name)}</h3>
                <div class="skill-card-badges">
                    ${sfiaBadge}
                    ${frameworkBadge}
                </div>
            </div>
            <p class="skill-card-desc">${escapeHtml(skill.description)}</p>
            <div class="skill-card-tags">
                ${categoryBadge}
                ${skill.tags.slice(0, 3).map(tag => `<span class="tag">${escapeHtml(tag)}</span>`).join('')}
            </div>
        </div>
    `;
}

async function openSkillModal(slug) {
    const modal = document.getElementById('skillModal');
    const modalBody = document.getElementById('modalBody');

    try {
        const res = await fetch(`${API_BASE}/skills/${slug}`);
        if (!res.ok) throw new Error('Skill not found');
        
        const data = await res.json();
        const skill = data.data;

        const sfiaBadge = skill.sfia ? 
            `<span class="badge sfia">SFIA Level ${skill.sfia.level}</span>` : '';
        const frameworkBadge = skill.framework ? 
            `<span class="badge framework">${skill.framework.type.toUpperCase()}</span>` : '';

        modalBody.innerHTML = `
            <div class="skill-detail">
                <div class="skill-detail-header">
                    <div>
                        <h2 class="skill-detail-title">${escapeHtml(skill.name)}</h2>
                        <p class="skill-detail-desc">${escapeHtml(skill.description)}</p>
                    </div>
                    <div class="skill-detail-badges">
                        ${sfiaBadge}
                        ${frameworkBadge}
                        ${skill.category ? `<span class="badge category">${escapeHtml(skill.category)}</span>` : ''}
                    </div>
                </div>
                
                <div class="skill-detail-content">
                    <aside class="skill-sidebar">
                        ${skill.tags.length ? `
                            <div class="sidebar-section">
                                <h5>Tags</h5>
                                <div class="skill-tags">
                                    ${skill.tags.map(tag => `<span class="tag">${escapeHtml(tag)}</span>`).join('')}
                                </div>
                            </div>
                        ` : ''}

                        ${skill.sfia ? `
                            <div class="sidebar-section">
                                <h5>SFIA Level</h5>
                                <p><strong>Level ${skill.sfia.level}:</strong> ${escapeHtml(skill.sfia.competency || 'Not specified')}</p>
                                ${skill.sfia.skills?.length ? `
                                    <p class="sfia-skills">${skill.sfia.skills.join(', ')}</p>
                                ` : ''}
                            </div>
                        ` : ''}

                        ${skill.qualityMetrics ? `
                            <div class="sidebar-section">
                                <h5>Quality Metrics</h5>
                                <div class="quality-bars">
                                    <div class="quality-bar">
                                        <span>Accuracy</span>
                                        <div class="bar-bg"><div class="bar-fill" style="width:${skill.qualityMetrics.accuracy * 100}%"></div></div>
                                        <span class="bar-value">${(skill.qualityMetrics.accuracy * 100).toFixed(0)}%</span>
                                    </div>
                                    <div class="quality-bar">
                                        <span>Consistency</span>
                                        <div class="bar-bg"><div class="bar-fill" style="width:${skill.qualityMetrics.consistency * 100}%"></div></div>
                                        <span class="bar-value">${(skill.qualityMetrics.consistency * 100).toFixed(0)}%</span>
                                    </div>
                                    <div class="quality-bar">
                                        <span>Completeness</span>
                                        <div class="bar-bg"><div class="bar-fill" style="width:${skill.qualityMetrics.completeness * 100}%"></div></div>
                                        <span class="bar-value">${(skill.qualityMetrics.completeness * 100).toFixed(0)}%</span>
                                    </div>
                                </div>
                            </div>
                        ` : ''}

                        ${skill.guardrails ? `
                            <div class="sidebar-section">
                                <h5>Guardrails</h5>
                                ${skill.guardrails.intendedUse?.length ? `<p><strong>Use:</strong> ${skill.guardrails.intendedUse.join(', ')}</p>` : ''}
                                ${skill.guardrails.constraints?.length ? `<p><strong>Constraints:</strong> ${skill.guardrails.constraints.join(', ')}</p>` : ''}
                            </div>
                        ` : ''}

                        <div class="sidebar-actions">
                            <button class="btn btn-primary btn-block" onclick="copyPrompt('${slug}')">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                    <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                                    <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
                                </svg>
                                Copy Full Prompt
                            </button>
                            <button class="btn btn-ghost btn-block" onclick="addToBundle('${slug}')">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                    <path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z"/>
                                    <line x1="3" y1="6" x2="21" y2="6"/>
                                    <path d="M16 10a4 4 0 0 1-8 0"/>
                                </svg>
                                Add to Bundle
                            </button>
                            <button class="btn btn-ghost btn-block" onclick="downloadSingleSkill('${slug}')">
                                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                                    <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
                                    <polyline points="7 10 12 15 17 10"/>
                                    <line x1="12" y1="15" x2="12" y2="3"/>
                                </svg>
                                Download .md
                            </button>
                        </div>
                    </aside>
                    
                    <div class="skill-prompt-viewer">
                        <div class="prompt-header">
                            <h5>Prompt Template</h5>
                            <span class="prompt-length">${skill.promptTemplate.length.toLocaleString()} characters</span>
                        </div>
                        <pre class="prompt-content">${escapeHtml(skill.promptTemplate)}</pre>
                    </div>
                </div>
            </div>
        `;

        modal.classList.add('open');
    } catch (err) {
        showToast('Failed to load skill details');
        console.error(err);
    }
}

function closeModalFn(modal) {
    modal.classList.remove('open');
}

async function copyPrompt(slug) {
    try {
        const res = await fetch(`${API_BASE}/skills/${slug}`);
        const data = await res.json();
        const skill = data.data;

        await navigator.clipboard.writeText(skill.promptTemplate);
        showToast('Prompt copied to clipboard');
    } catch (err) {
        showToast('Failed to copy prompt');
        console.error(err);
    }
}

function addToBundle(slug) {
    if (state.bundle.includes(slug)) {
        showToast('Skill already in bundle');
        return;
    }

    state.bundle.push(slug);
    updateBundleUI();
    showToast('Added to bundle');
}

function removeFromBundle(slug) {
    state.bundle = state.bundle.filter(s => s !== slug);
    updateBundleUI();
}

function updateBundleUI() {
    const badge = document.getElementById('bundleBadge');
    const toggle = document.getElementById('bundleToggle');
    const count = document.getElementById('bundleCount');
    const items = document.getElementById('bundleItems');
    const downloadBtn = document.getElementById('downloadZip');

    const bundleCount = state.bundle.length;
    badge.textContent = bundleCount;
    count.textContent = `${bundleCount} item${bundleCount !== 1 ? 's' : ''}`;
    downloadBtn.disabled = bundleCount === 0;

    toggle.style.display = bundleCount > 0 ? 'flex' : 'none';

    items.innerHTML = state.bundle.map(slug => {
        const skill = state.skills.find(s => s.slug === slug);
        return `
            <div class="bundle-item">
                <span class="bundle-item-name">${escapeHtml(skill?.name || slug)}</span>
                <button class="bundle-item-remove" onclick="removeFromBundle('${slug}')">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="16" height="16">
                        <path d="M18 6L6 18M6 6l12 12"/>
                    </svg>
                </button>
            </div>
        `;
    }).join('');
}

function toggleBundle(open) {
    const drawer = document.getElementById('bundleDrawer');
    const isOpen = drawer.classList.contains('open');
    
    if (open === undefined) {
        drawer.classList.toggle('open', !isOpen);
    } else {
        drawer.classList.toggle('open', open);
    }
}

async function downloadBundleZip() {
    if (state.bundle.length === 0) {
        showToast('No items in bundle');
        return;
    }

    showToast('Preparing bundle...');

    try {
        const res = await fetch(`${API_BASE}/bundle`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ slugs: state.bundle })
        });

        if (!res.ok) {
            throw new Error('Failed to fetch bundle');
        }

        const data = await res.json();
        const skills = data.data;

        if (!skills || skills.length === 0) {
            showToast('No skills found for bundle');
            console.error('Bundle API returned empty:', data);
            return;
        }

        const manifest = {
            version: '1.0',
            createdAt: new Date().toISOString(),
            count: skills.length,
            skills: skills.map(s => ({
                slug: s.slug,
                name: s.name,
                version: s.version
            }))
        };

        const content = `# Promptito Bundle
Generated: ${new Date().toLocaleString()}
Total Items: ${skills.length}

---

${skills.map(s => `# ${s.name}

**Slug:** ${s.slug}
**Category:** ${s.category || 'N/A'}
**Tags:** ${(s.tags || []).join(', ')}

---

${s.promptTemplate}

---
`).join('\n')}`;

        downloadFile(`promptito-bundle-${Date.now()}.md`, content, 'text/markdown');
        showToast(`Downloaded ${skills.length} prompt${skills.length > 1 ? 's' : ''}`);
    } catch (err) {
        showToast('Failed to download bundle');
        console.error('Bundle download error:', err);
    }
}

function downloadFile(filename, content, type) {
    const blob = new Blob([content], { type });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
}

async function downloadSingleSkill(slug) {
    try {
        const res = await fetch(`${API_BASE}/skills/${slug}`);
        const data = await res.json();
        const skill = data.data;

        const yaml = `---
name: ${skill.name}
version: ${skill.version}
description: ${skill.description}
${skill.category ? `category: ${skill.category}` : ''}
${skill.tags?.length ? `tags:\n${skill.tags.map(t => `  - ${t}`).join('\n')}` : ''}
${skill.sfia ? `sfia:\n  level: ${skill.sfia.level}\n  competency: ${skill.sfia.competency || ''}` : ''}
${skill.qualityMetrics ? `qualityMetrics:\n  accuracy: ${skill.qualityMetrics.accuracy}\n  consistency: ${skill.qualityMetrics.consistency}\n  completeness: ${skill.qualityMetrics.completeness}` : ''}
---
${skill.promptTemplate}`;

        downloadFile(`${slug}.md`, yaml, 'text/markdown');
        showToast('Downloaded ' + skill.name);
    } catch (err) {
        showToast('Failed to download');
    }
}

function showToast(message) {
    const toast = document.getElementById('toast');
    toast.textContent = message;
    toast.classList.add('show');
    setTimeout(() => toast.classList.remove('show'), 3000);
}

function debounce(fn, delay) {
    let timeout;
    return (...args) => {
        clearTimeout(timeout);
        timeout = setTimeout(() => fn(...args), delay);
    };
}

function escapeHtml(str) {
    if (!str) return '';
    const div = document.createElement('div');
    div.textContent = str;
    return div.innerHTML;
}
