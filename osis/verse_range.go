package osis

//GetVerses takes two reference strings and returns a slice of verses
func (b *Bible) GetVerses(from, to string) ([]*Verse, error) {
	if b.Books == nil || len(b.Books) == 0 {
		b.index()
	}

	fromRef, err := NewVerseRef(from)
	if err != nil {
		return nil, err
	}
	toRef, err := NewVerseRef(to)
	if err != nil {
		return nil, err
	}

	var firstBook, lastBook int
	var ok bool
	if firstBook, ok = b.BooksById[fromRef.BookID]; !ok {
		return nil, NoSuchVerse
	}
	if lastBook, ok = b.BooksById[toRef.BookID]; !ok {
		return nil, NoSuchVerse
	}

	verses := make([]*Verse, 0, 100)
	for currBook := firstBook; currBook <= lastBook; currBook++ {
		book := b.Books[currBook]

		var firstChap, lastChap int
		if currBook == firstBook {
			firstChap = fromRef.Chapter
		}
		if currBook == lastBook {
			lastChap = toRef.Chapter
		} else {
			lastChap = len(book.Chs) - 1
		}

		for currChap := firstChap; currChap <= lastChap; currChap++ {
			ch := book.Chs[currChap]

			var firstVerse, lastVerse int
			if currChap == firstChap {
				firstVerse = fromRef.Verse
			}

			if currChap == lastChap {
				lastVerse = toRef.Verse
			} else {
				lastVerse = len(ch.Vrs)
			}

			verses = append(verses, ch.Vrs[firstVerse:lastVerse+1]...)
		}
	}

	return verses, nil
}
