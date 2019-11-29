from pathlib import Path # recoursive walking
import operator # sorted dicts
import re # The dvornik is dead. Long live the dvornik
import json

def dvornik(article: str) -> str:
    """ Delete all html tags
    """

    # TODO: delete all trash such as !?,.\n etc in this func
    # Becouse regexp is sucks

    length = len(article)
    memory_carrage = 0
    carrage = 0

    while True:
        if article[carrage] == "<":
            memory_carrage = carrage
        if article[carrage] == ">" and memory_carrage is not None:
            article = article[:memory_carrage] + article[carrage+1:]
            length -= carrage - memory_carrage + 1
            carrage -= carrage - memory_carrage + 1
            memory_carrage = None

        carrage += 1
        if carrage == length:
            break

    # XXX
    # TODO: Это некий костыль, я возвращаю обрезанную часть
    # текста из этой функции. Оформить в виде константы, чтобы
    # можно было менять. Очень сильно влияет на скорость обработки
    # и на качество предсказания. Перенести вообще в другую функцию,
    # например в detect_language
    return article

def bi_grams(sentence: str) -> str:
    """ Construct array of bi-grams, sorted with respect
    on frequency. Actually, phonemes != bi-grams
    """
    sentence = re.sub(r"[,.!-?]+", "", sentence.lower())

    phonemes = dict()
    words = sentence.split(" ")

    for word in words:
        length = len(word)
        for ind,char in enumerate(word):
            if (n:=ind+1) < length:
                if (ph:=f"{char}{word[n]}") in phonemes:
                    phonemes[ph] += 1
                else:
                    phonemes[ph] = 1

    sorted_phonemes_hash_table = sorted(
        phonemes.items(),
        key=operator.itemgetter(1),
        reverse=True
    )

    return list(map(lambda x: x[0], sorted_phonemes_hash_table))

def out_of_place_measure(phonemes, corpus_phonemes) -> int:
    """ The measure between two sorted array of bi-grams.
    The size of two arrays must be equal. The result will be
    equal len(phonemes)^2 in the worst-case scenario.
    """
    distance = 0
    for ph in phonemes:
        if ph not in corpus_phonemes:
            distance += len(corpus_phonemes)
        else:
            distance += abs(
                corpus_phonemes.index(ph)-phonemes.index(ph)
            )
    return distance

def load_lang_profiles(file_path="language_profiles.json") -> dict:
    with open("language_profiles.json", "r") as file:
        return json.load(file)

def detect_language(data, amount=None) -> str:
    """ Predict language of the text using bi_grams method
    and out-Of-Place measure
    """
    language_profiles = load_lang_profiles()
    
    data_grams = bi_grams(data)
    

    # Calculate amount of bi_grams
    # amount <= 445
    if len(data_grams) > 445:
        data_grams = data_grams[:445]
    if amount is None:
        amount = len(data_grams)
    else:
        if amount > 445:
            amount = len(data_grams)

    measure = amount ** 2
    predicted_lang = "other"

    for lang in language_profiles.keys():
        if (curr_measure := out_of_place_measure(language_profiles[lang][:amount], data_grams[:amount])) < measure:
            measure = curr_measure
            predicted_lang = lang
    return predicted_lang


if __name__ == "__main__":
    # example of using

    language_profiles = load_lang_profiles()

    with open(file_path, "r") as file:
        data = f.read()
    data = dvornik(data)
    l = detect_language(data)
    print(l)
