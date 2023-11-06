package main

import (
    "fmt"
    "image/color"
    "image/jpeg" // Import the JPEG package
    "log"
    "os"
    "path/filepath"
    "github.com/fogleman/gg"
)

func main() {
    // Define the watermark text and its properties
    watermarkText := "Copythight (C) Amila Hewagama"
    fontSize := 48.0

    // Directory containing images
    imageDir := "images"

    // Create the "output" directory if it doesn't exist
    outputDir := "output"
    if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
        log.Fatal(err)
    }

    // Open the directory and loop through its contents
    err := filepath.Walk(imageDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        // Check if the file is an image (you may want to add more image file extensions)
        if !info.IsDir() && (filepath.Ext(path) == ".jpg" || filepath.Ext(path) == ".jpeg" || filepath.Ext(path) == ".png") {
            // Open the image
            inputImage, err := gg.LoadImage(path)
            if err != nil {
                return err
            }

            // Create a context for drawing on the image
            dc := gg.NewContextForImage(inputImage)

            // Set font properties
            fontFamily := "Arial" // Change this to the desired font family
            dc.LoadFontFace(fontFamily, fontSize)

            // Calculate text width and height
            textWidth, textHeight := dc.MeasureString(watermarkText)

            // Calculate the position to center the watermark
            imageWidth := float64(inputImage.Bounds().Dx())
            imageHeight := float64(inputImage.Bounds().Dy())
            textX := (imageWidth - textWidth) / 2
            textY := (imageHeight - textHeight) / 2

            // Set the text color with transparency
            dc.SetColor(color.RGBA{0, 0, 0, 128}) // 128 for 50% opacity (adjust as needed)
            dc.DrawStringAnchored(watermarkText, textX, textY, 0.5, 0.5) // Adjust the last two arguments for alignment

            // Create an output file in the "output" directory
            outputPath := filepath.Join(outputDir, info.Name())
            outputFile, err := os.Create(outputPath)
            if err != nil {
                return err
            }
            defer outputFile.Close()

            // Encode and save the result as JPEG with the same quality
            jpegOptions := jpeg.Options{Quality: 100} // Use 100 for the same quality level
            err = jpeg.Encode(outputFile, dc.Image(), &jpegOptions)
            if err != nil {
                return err
            }

            // Print a success message
            fmt.Printf("Watermark added to %s and saved as %s\n", path, outputPath)
        }

        return nil
    })

    if err != nil {
        log.Fatal(err)
    }
}
