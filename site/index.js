document.addEventListener('DOMContentLoaded', () => {
  const input = document.querySelector('.neo-input');
  const validateBtn = document.querySelector('.btn.validate');
  const confirmBtn = document.querySelector('.btn.confirm');
  const slider = document.getElementById('lengthSlider');
  const valueDisplay = document.getElementById('lengthValue');

  valueDisplay.textContent = slider.value;

  slider.addEventListener('input', () => {
    valueDisplay.textContent = slider.value;
  });
  const showResult = (resultText) => {
   
    let resultEl = document.getElementById('result');
    if (!resultEl) {
      resultEl = document.createElement('div');
      resultEl.id = 'result';
      resultEl.style.marginTop = '1rem';
      resultEl.style.fontFamily = 'sans-serif';
      resultEl.style.fontWeight = 'bold';
      document.querySelector('.neo-form').appendChild(resultEl);
    }
    resultEl.textContent = resultText;
  };

  validateBtn.addEventListener('click', async (e) => {
    e.preventDefault();
    const value = input.value;
    try {
      const resp = await fetch('/check', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ url: value }),
      });
      if (!resp.ok) throw new Error(`HTTP ${resp.status}`);
      let text = await resp.text();
      if(!text) {
        text = "Could not be found"
      }
      showResult(`POST result: ${text}`);
    } catch (err) {
      showResult(`POST error: ${err}`);
    }
  });

  confirmBtn.addEventListener('click', async (e) => {
    e.preventDefault();
    const value = input.value;
    const length = slider.value;
    try {
      const resp = await fetch('/create', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ url: value , length}),
      });
      if (!resp.ok) throw new Error(`HTTP ${resp.status}`);
      const text = await resp.text();
      showResult(`POST result: ${text}`);
    } catch (err) {
      showResult(`POST error: ${err}`);
    }
  });
});
