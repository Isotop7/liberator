// Attach to slider
bulmaSlider.attach()

slider = document.getElementById("sliderWithValue");
slider.addEventListener("input", () => {
    sliderOutput = document.getElementById("sliderOutput");
    val = parseInt(slider.value)
    if (val < 4) {
        sliderOutput.className = "has-text-black has-background-danger"
    } else if (val < 8) {
        sliderOutput.className = "has-text-black has-background-warning"
    } else {
        sliderOutput.className = "has-text-black has-background-success"
    }
});