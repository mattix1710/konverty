from ffmpeg_quality_metrics import FfmpegQualityMetrics
import json
from pathlib import Path

distorted = []

PATH_METRICS = Path(Path.cwd(), "metrics")

for el in Path(Path.cwd(), "examples/dist").iterdir():
    if el.is_file():
        distorted.append([el, el.stem])

for el in distorted:
    ffqm = FfmpegQualityMetrics(str(Path(Path.cwd(), "examples/ref/input.mp4")), str(el[0]))
    
    metrics = ffqm.calculate(["psnr", "ssim"])
    
    with open(Path(PATH_METRICS, "file_{}_psnr.json".format(el[1])), "w") as file_json:
        json.dump(metrics["psnr"], file_json, indent=4)
    with open(Path(PATH_METRICS, "file_{}_ssim.json".format(el[1])), "w") as file_json:
        json.dump(metrics["ssim"], file_json, indent=4)
    
    print("INFO: File {} processed.".format(str(el[0])))