TARGET := usenix23

TEX_ALL := $(shell search=$(TARGET).tex; all=; \
				while [ -n "$$search" ] ; do \
					all="$$all $$search"; \
					search=`egrep "^[^%]*\\input" $$search | \
						sed -En 's/.*input[^\{]*\{(.+)\}/\1.tex/p'`; \
				done; \
				echo "$$all")

FIGURES :=

all: $(TARGET).pdf

figs/%.eps: plots/%.plot
	gnuplot $< > $@

$(TARGET).pdf: $(FIGURES)  $(TEX_ALL)
	pdflatex $(TARGET).tex
	# bibtex $(TARGET)
	pdflatex $(TARGET).tex
	pdflatex $(TARGET).tex

nobib:: $(addsuffix .eps, $(FIGURES)) $(TARGET).aux
	dvips $(TARGET).dvi -Ppdf -Pcmz -Pamz -t letter -D 600 -G0 -o $(TARGET).ps
	ps2pdf ${PS2PDF_FLAGS} $(TARGET).ps

view:: $(TARGET).pdf
	open $(TARGET).pdf

see:: $(TARGET).dvi
	xdvi $(TARGET)

%.eps: plots/%.plot
	gnuplot $< > figs/$@

%.eps: figs/%.fig
	fig2dev -L eps -p dummy $< $*.eps

spell::
	ispell *.tex

clean::
	rm -fv *.dvi *.aux *.log *~ *.bbl *.blg *.toc *.out *.ps *.ent parsetab.py
	rm -fv $(addsuffix -eps-converted-to.pdf, $(FIGURES))

fresh::
	rm -fv *.dvi *.aux *.log *~ *.bbl *.blg *.toc *.ps *.pdf

distclean:: clean
	rm $(TARGET).ps

.PHONY: clean all view see spell fresh distclean nobib
