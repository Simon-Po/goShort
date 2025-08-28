document.addEventListener('DOMContentLoaded', () => {
  const form = document.querySelector('.neo-form');
  const urlInput = document.querySelector('#user_input');
  const nameInput = document.querySelector('#urlName');
  const slider = document.getElementById('lengthSlider');
  const sliderLabel = document.querySelector('label[for="lengthSlider"]');
  const validateBtn = document.querySelector('.btn.validate');
  const confirmBtn = document.querySelector('.btn.confirm');

  const showResult = (t) => {
    let el = document.getElementById('result');
    if (!el) { el = document.createElement('div'); el.id='result'; form.appendChild(el); }
    el.textContent = t;
  };

  // 👇 toggle slider visibility based on urlName input
  nameInput.addEventListener('input', () => {
    if (nameInput.value.trim() !== '') {
      slider.style.display = 'none';
      sliderLabel.style.display = 'none';
    } else {
      slider.style.display = 'block';
      sliderLabel.style.display = 'block';
    }
  });

  confirmBtn.addEventListener('click', async (e) => {
    e.preventDefault();
    if (!form.reportValidity()) return;

    const payload = {
      name: nameInput.value,
      url: urlInput.value,
      length: String(slider.value)
    };

    console.log('about to POST /create with payload:', payload);

    try {
      const resp = await fetch('/create', {
        method: 'POST',
        headers: {'Content-Type':'application/json'},
        body: JSON.stringify(payload),
      });
      const text = await resp.text();
      showResult(`POST result: ${text}`);
    } catch (err) {
      showResult(`POST error: ${err}`);
    }
  });

  validateBtn.addEventListener('click', async (e) => {
    e.preventDefault();
    try {
      const resp = await fetch('/check', {
        method: 'POST',
        headers: {'Content-Type':'application/json'},
        body: JSON.stringify({ url: urlInput.value }),
      });
      const text = await resp.text();
      showResult(`POST result: ${text || 'Could not be found'}`);
    } catch (err) {
      showResult(`POST error: ${err}`);
    }
  });
});
