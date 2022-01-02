# Creating GIFs for Tidbyt

By using Pixlet's [image widget](widgets.md), images and GIFs can easily be displayed on pixel constrained displays. 

However, when creating GIFs for use on Tidbyt, there are a few requirements to keep in mind for the best result.

## Design GIFs to be 64x32 pixels from the start
Tidbyt's display is 64x32 pixels. If there is a GIF that’s larger than 64x32 pixels, it has to be scaled down. In practice, we’ve found that images scaled down to this resolution don’t look as crisp as when images are designed for 64x32 from the beginning. So if you’re creating GIFs for the Tidbyt, make sure they’re 64x32.

## Finished GIF is 128KB or less
We’re limited by the number of bytes we can send to the Tidbyt and the Tidbyt is constrained by how many bytes it can store locally. To get around this, we limit the size of the GIF to 128 Kilobytes and if it’s larger than this after downsizing to 64x32 pixels, we drop frames until it fits the size requirements. This means if you want your GIF to look great on the Tidbyt, make sure it’s 128KB or less before adding it through the mobile app or with Pixlet.

## GIF is 15 seconds in length or loops cleanly if less then 15 seconds.
The length of time the GIF loops should be around 15 seconds. The timings for applet cycles are 15, 10, 7.5, and 5 seconds depending on the setting in the mobile app. If your GIF is less then 15 seconds, ensure it loops cleanly to avoid an interrupt.