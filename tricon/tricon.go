/*
  tricon.go
  SDL Debug Console
  version: 17.12.23
  Copyright (C) 2017 Jeroen P. Broks
  This software is provided 'as-is', without any express or implied
  warranty.  In no event will the authors be held liable for any damages
  arising from the use of this software.
  Permission is granted to anyone to use this software for any purpose,
  including commercial applications, and to alter it and redistribute it
  freely, subject to the following restrictions:
  1. The origin of this software must not be misrepresented; you must not
     claim that you wrote the original software. If you use this software
     in a product, an acknowledgment in the product documentation would be
     appreciated but is not required.
  2. Altered source versions must be plainly marked as such, and must not be
     misrepresented as being the original software.
  3. This notice may not be removed or altered from any source distribution.
*/

// Tricon is just a very simplistic graphic debug console.
// It uses veandco's SDL library, also meaning that on Mac it can only be used fully set up Application bundles.
package tricon

import (
	"trickyunits/mkl"
	"trickyunits/qstr"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)


type pieceoftext struct {
	x int32
	y int32
	txt string
	r uint8
	g uint8
	b uint8
	s *sdl.Surface
	t *sdl.Texture
}

func (rd *pieceoftext) kill(){
	rd.s.Free()
	rd.t.Destroy()
}

var rend *sdl.Renderer
var font *ttf.Font
var x int32
var y int32
var ww int32
var wh int32
var txl []pieceoftext

var white = sdl.Color{R:255,G:255,B:255,A:255}



// When set the console will output all data into an output file.
// The output will be in html in order keep color settings as well as possible in order, although the system is not sophisticated enough to copy stuff exactly the way they are put on the screen.
// If you want a clean setting it would also be wise to delete this file (if it exists) at the start of your session, or the date of earlier sessions will linger there.
// The file will be created on the moment the first data is outputted to the console. It would be wise to have color settings set before that happens!
var Outfile string

// Used for the background color. 
var BR,BG,BB uint8

// Used for the color of the command line
var CR,CG,CB uint8 = 255,255,255

var cmd string


// Last Command that has been confirmed by the user
var LastCommand string


// Must be created with a Window as example or tricon never knows how big its workspace is.
func Setup(w *sdl.Window,f *ttf.Font){
	//var err error
	rend,_ = sdl.CreateRenderer(w,-1,0)
	font = f
	txl = []pieceoftext{}
	ww,wh = w.GetSize()
}

// Shows the current state of the console.
func Show() {
	//var err error
	rend.SetDrawColor(BR,BG,BB,255)
	rend.Clear()
	//updated:=[]pieceoftext{}
	for _,tx:=range txl{
		if tx.s==nil {
			tx.s,_ = font.RenderUTF8Blended(tx.txt,white)
			tx.t,_ = rend.CreateTextureFromSurface(tx.s)
		}
		tx.t.SetColorMod(tx.r,tx.g,tx.b)
		//rend.SetDrawColor(tx.r,tx.g,tx.b,255)		
		src := sdl.Rect{0, 0, tx.s.W, tx.s.H}
		dst := sdl.Rect{tx.x, tx.y, tx.s.W, tx.s.H}	
		rend.Copy(tx.t,&src,&dst)
	}
	csurf,_:=font.RenderUTF8Blended(">"+cmd+"_",white); defer csurf.Free()
	ctext,_:=rend.CreateTextureFromSurface(csurf); defer ctext.Destroy()
	ctext.SetColorMod(CR,CG,CB)
	ctsr:=sdl.Rect{0, 0, csurf.W, csurf.H}
	cttr:=sdl.Rect{0,int32(wh)-int32(font.Height()),csurf.W,csurf.H}
	rend.Copy(ctext,&ctsr,&cttr)
	rend.Present()
}

func Write(txt string,r,g,b uint8){
	i:=pieceoftext{}
	i.x=x
	i.y=y
	i.r=r
	i.g=g
	i.b=b
	i.txt = txt
	i.s,_ = font.RenderUTF8Blended(txt,white)
	i.t,_ = rend.CreateTextureFromSurface(i.s)
	x+=i.s.W
	txl = append(txl,i)
}

func WriteLn(txt string,r,g,b uint8){
	Write(txt,r,g,b)
	fh:=int32(font.Height())
	x=0
	y+=fh
	if y>wh-(fh) {
		//updated:=[]pieceoftext{}
		for i:=0;i<len(txl);i++ {
			txl[i].y -= fh
			/*
			if txl[i].y < 0-(fh*2) {
				txl[i].kill()
			} else {
				updated=append(updated,txl[i])
			}
			*/
		//txl=updated
		}
		y-=fh
		for txl[0].y< 0-(fh+fh){
			txl[0].kill()
			txl = txl[1:]
		}
	}
}

// Adds extra text to the current command on the bottom.
func CADD(t string) { cmd+=t }

// Resets the command bar and returs its content
func Confirm() string {
	ret:=cmd
	WriteLn(">"+cmd,CR,CG,CB)
	cmd=""
	return ret
}

// Backspace
func CBS() {
	if cmd=="" { return }
	cmd = qstr.Left(cmd,len(cmd)-1)
	
}

// Reset
func CRS() { cmd="" }


// A traditional function in my own right. Don't pay it any mind :P
func CSay(txt string){
	WriteLn(txt,255,180,0)
}

// Frees everything taking up memory so far!
func Kill(){
	for i:=0;i<len(txl);i++ {
		txl[i].kill()
	}
	rend.Destroy()
}

func init(){
mkl.Version("Tricky's Go Units - tricon.go","17.12.23")
mkl.Lic    ("Tricky's Go Units - tricon.go","ZLib License")
}
