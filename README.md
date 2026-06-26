# Konverty

1. FFmpeg convertion:
`ffmpeg -i <input> -c:v libx265 -crf 23 -preset fast -c:a aac -b:a 192k <output>`
2. Probing for bitrate value:
`ffprobe -v error -select_streams a:0 -show_entries stream=bit_rate -of default=noprint_wrappers=1:nokey=1 <input_file>`
3. Probing for codec name:
`ffprobe -v error -select_streams a:0 -show_entries stream=codec_name -of default=noprint_wrappers=1:nokey=1 <input_file>`
4. Probing for frame count:
`ffprobe -v error -select_streams v:0 -count_packets -show_entries stream=nb_read_packets -of csv=p=0 <input_file>`
5. Calculating PSNR:
`ffmpeg -i <original_file> -i <processed_file> -v info -stats -lavfi psnr -f null -`

Misc cmds:

1. Cut out:
`ffmpeg -i '.\Lata 2003-04 (2).avi' -ss 00:08:27 -to 00:08:35 -c copy z_wody.mp4`

- With reencoding:
`ffmpeg -i '.\Rok 2002.avi' -ss 00:21:14 -to 00:21:22 -c:v libx265 -c:a aac mini-ja.mp4`

## TODO

- [ ] Create channel for keyboard interruption - Shutting down gracefully when converting and deleting partially created video
- [x] ~~Add logging for size differences in conversion~~
- [ ] Add video file detection (IDEA: based on video codec format detection?)
- [ ] Lower PSNR threshold (curr: 30)
- [ ] Add flags based on mode type (eg. `-q` for quality assessment, `-X` for overwriting without prompting)
