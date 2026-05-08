from PIL import Image, ImageDraw
import os

def create_icon():
    size = 256
    img = Image.new('RGBA', (size, size), (0, 0, 0, 0))
    draw = ImageDraw.Draw(img)

    green = (29, 185, 84, 255)

    draw.ellipse([55, 140, 105, 185], fill=green)
    draw.rectangle([100, 50, 112, 165], fill=green)
    draw.polygon([(112, 50), (112, 100), (155, 75)], fill=green)

    draw.ellipse([140, 100, 175, 130], fill=green)
    draw.rectangle([168, 30, 178, 115], fill=green)

    return img

def create_ico(png_path, ico_path):
    img = Image.open(png_path).convert('RGBA')
    sizes = [16, 32, 48, 64, 128, 256]
    imgs = [img.resize((s, s), Image.Resampling.LANCZOS) for s in sizes]
    imgs[0].save(ico_path, format='ICO', sizes=[(s, s) for s in sizes], append_images=imgs[1:])

def main():
    base = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    build_dir = os.path.join(base, 'build')
    windows_dir = os.path.join(build_dir, 'windows')

    os.makedirs(windows_dir, exist_ok=True)

    img = create_icon()

    png_path = os.path.join(build_dir, 'icon_256.png')
    img.save(png_path, 'PNG')
    print(f"[OK] {png_path}")

    ico_path = os.path.join(windows_dir, 'icon_tray.ico')
    create_ico(png_path, ico_path)
    print(f"[OK] {ico_path}")

    notif_path = os.path.join(windows_dir, 'icon_notification.png')
    img.resize((128, 128), Image.Resampling.LANCZOS).save(notif_path, 'PNG')
    print(f"[OK] {notif_path}")

if __name__ == '__main__':
    main()