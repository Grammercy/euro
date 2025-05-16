document.addEventListener('DOMContentLoaded', () => {
    loadFigures();
    setupModal();
    setupTimelineObserver();
    setupRankingSystem();
});

let figures = []; // Store figures globally

async function loadFigures() {
    try {
        const response = await fetch('/api/figures');
        figures = await response.json();
        
        // Populate the ranking container
        const rankingContainer = document.getElementById('ranking-container');
        figures.forEach(figure => {
            const card = createFigureCard(figure);
            card.classList.add('ranking-card');
            card.dataset.figureId = figure.id;
            rankingContainer.appendChild(card);
        });

        // Create timeline
        createTimeline(figures);
    } catch (error) {
        console.error('Error loading figures:', error);
    }
}

function createFigureCard(figure) {
    const card = document.createElement('div');
    card.className = 'figure-card bg-white rounded-lg shadow-lg overflow-hidden relative';
    card.innerHTML = `
        <div class="rank-badge">${figure.rank || '?'}</div>
        <img src="${figure.imageUrl}" alt="${figure.name}" class="w-full h-48 object-cover">
        <div class="p-4">
            <h3 class="text-xl font-bold">${figure.name}</h3>
            <p class="text-gray-600">${figure.period}</p>
            <p class="mt-2">${figure.innovation}</p>
        </div>
    `;
    
    card.addEventListener('click', () => showModal(figure));
    return card;
}

function createTimeline(figures) {
    const timeline = document.querySelector('.timeline-container');
    timeline.innerHTML = ''; // Clear existing content
    
    figures.forEach(figure => {
        const item = document.createElement('div');
        item.className = 'timeline-item';
        
        // Extract years from period
        const years = figure.period.split('-').map(year => year.trim());
        
        item.innerHTML = `
            <div class="timeline-content">
                <div class="timeline-year">${figure.period}</div>
                <h4 class="timeline-title">${figure.name}</h4>
                <p class="text-sm text-gray-600">${figure.innovation}</p>
            </div>
        `;
        
        // Add click event to timeline items
        item.addEventListener('click', () => showModal(figure));
        timeline.appendChild(item);
    });
}

function setupTimelineObserver() {
    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                entry.target.classList.add('visible');
            }
        });
    }, {
        threshold: 0.2,
        rootMargin: '0px'
    });

    // Observe all timeline items
    document.querySelectorAll('.timeline-item').forEach(item => {
        observer.observe(item);
    });
}

function setupRankingSystem() {
    const rankingContainer = document.getElementById('ranking-container');
    const submitButton = document.getElementById('submit-ranking');
    const showAveragesButton = document.getElementById('show-averages');
    const averagesModal = document.getElementById('averages-modal');
    const closeAveragesBtn = document.querySelector('.close-averages');

    // Function to update rank badges
    function updateRankBadges() {
        const cards = rankingContainer.children;
        for (let i = 0; i < cards.length; i++) {
            const badge = cards[i].querySelector('.rank-badge');
            if (badge) {
                badge.textContent = (i + 1).toString();
            }
        }
    }

    // Initialize Sortable
    new Sortable(rankingContainer, {
        animation: 150,
        ghostClass: 'bg-blue-50',
        onEnd: function(evt) {
            updateRankBadges();
        }
    });

    // Initial rank badge update
    updateRankBadges();

    // Submit rankings
    submitButton.addEventListener('click', async () => {
        const rankings = Array.from(rankingContainer.children).map((card, index) => ({
            figureId: parseInt(card.dataset.figureId),
            rank: index + 1
        }));

        try {
            const response = await fetch('/api/rankings', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(rankings)
            });

            if (response.ok) {
                alert('Your rankings have been submitted successfully!');
            } else {
                alert('Error submitting rankings. Please try again.');
            }
        } catch (error) {
            console.error('Error submitting rankings:', error);
            alert('Error submitting rankings. Please try again.');
        }
    });

    // Show average rankings
    showAveragesButton.addEventListener('click', async () => {
        try {
            const response = await fetch('/api/rankings/average');
            const averages = await response.json();

            const averagesContent = document.getElementById('averages-content');
            averagesContent.innerHTML = averages.map((avg, index) => `
                <div class="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
                    <div class="flex items-center space-x-4">
                        <span class="text-2xl font-bold text-blue-600">${index + 1}</span>
                        <div>
                            <h4 class="font-bold">${avg.name}</h4>
                            <p class="text-sm text-gray-600">Average Rank: ${avg.average.toFixed(1)}</p>
                            <p class="text-xs text-gray-500">Last updated: ${new Date(avg.updatedAt).toLocaleString()}</p>
                        </div>
                    </div>
                </div>
            `).join('');

            averagesModal.classList.remove('hidden');
        } catch (error) {
            console.error('Error fetching average rankings:', error);
            alert('Error loading average rankings. Please try again.');
        }
    });

    // Close averages modal
    closeAveragesBtn.addEventListener('click', () => {
        averagesModal.classList.add('hidden');
    });

    averagesModal.addEventListener('click', (e) => {
        if (e.target === averagesModal) {
            averagesModal.classList.add('hidden');
        }
    });
}

function setupModal() {
    const modal = document.getElementById('modal');
    const closeBtn = modal.querySelector('.close-modal');
    
    closeBtn.addEventListener('click', () => {
        modal.classList.add('hidden');
    });
    
    modal.addEventListener('click', (e) => {
        if (e.target === modal) {
            modal.classList.add('hidden');
        }
    });
}

function showModal(figure) {
    const modal = document.getElementById('modal');
    const title = modal.querySelector('.modal-title');
    const body = modal.querySelector('.modal-body');
    
    title.textContent = figure.name;
    body.innerHTML = `
        <div class="flex flex-col md:flex-row gap-6">
            <img src="${figure.imageUrl}" alt="${figure.name}" class="w-full md:w-1/3 object-cover rounded-lg">
            <div class="flex-1">
                <p class="text-gray-600 mb-2">${figure.period}</p>
                <h4 class="font-bold mb-2">Key Innovation:</h4>
                <p class="mb-4">${figure.innovation}</p>
                <h4 class="font-bold mb-2">Description:</h4>
                <p>${figure.description}</p>
            </div>
        </div>
    `;
    
    modal.classList.remove('hidden');
} 