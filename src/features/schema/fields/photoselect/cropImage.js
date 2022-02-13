export default function getCroppedImg(imageSrc, pixelCrop) {
    const image = new Image();
    image.src = imageSrc;

    const canvas = document.createElement('canvas');
    canvas.width = 64;
    canvas.height = 32;

    const ctx = canvas.getContext('2d');
    ctx.drawImage(
        image,
        pixelCrop.x,
        pixelCrop.y,
        pixelCrop.width,
        pixelCrop.height,
        0,
        0,
        64,
        32
    );

    return canvas.toDataURL('image/jpeg').replace('data:image/jpeg;base64,', '');
}
