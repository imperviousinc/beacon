import os

def extend(pkg_dmg, packaging_dir, icon_file, dsstore_file):
    pkg_dmg.extend([
        '--icon', os.path.join(packaging_dir, icon_file),
        '--mkdir', '.background',
        '--copy',
            '{}/beacon_dmg_background.tiff:/.background/background.tiff'.format(
                packaging_dir),
        '--copy', '{}/{}:/.DS_Store'.format(packaging_dir, dsstore_file),
    ])
