from bs4 import BeautifulSoup as bs
import sys
from utils import get_html
import json
import time
import re

terms = ["SUMMER TERM", "SEMESTER ONE", "SEMESTER TWO", "SEMESTER THREE"]
utilterms = ["U1", "T1", "T2", "T3"]
def main():

    root = str("http://timetable.unsw.edu.au/" + str(sys.argv[1]) + "/")
    classutil_url = "http://classutil.unsw.edu.au/{0}_{1}.html"
    handbookroot = "https://www.handbook.unsw.edu.au/{0}/courses/" + str(sys.argv[1]) + "/"
    soup = bs(get_html(root + "subjectSearch.html"), 'html.parser')

    courseinfo = {}
    courses = 0
    # for each subject area
    for area in soup.find_all('a'):
        arealink = area.get('href')
        if arealink is not None  and ".html" in arealink and len(arealink) <  14:
            groupsoup = bs(get_html(root+arealink), 'html.parser')
            currarea = arealink[:-5]
            print("GROUP: " + currarea)
            courseinfo[currarea] = {}
            if groupsoup is not None:
                for course in groupsoup.find_all('a'):
                    if course is not None and len(course) < 14:
                        courselink = course.get('href')
                        if courselink is not None and ".html" in courselink and len(courselink) < 14:                            
                            course_code = courselink[:-5]
                            print(course_code)
                            courseinfo[currarea][course_code] = {}
                            coursetimetable = str(get_html(root+courselink))                            
                            

                            # Semesters
                            courseinfo[currarea][course_code]["available_sems"] = []
                            for i in range(len(terms)):
                                if terms[i] in coursetimetable:
                                    courseinfo[currarea][course_code]["available_sems"].append(i)

                            # levels
                            undergradcourse = False
                            courseinfo[currarea][course_code]["levels"] = []
                            if "Undergraduate" in coursetimetable:
                                undergradcourse = True
                                courseinfo[currarea][course_code]["levels"].append("UGRAD")
                            if "Postgraduate" in coursetimetable:
                                courseinfo[currarea][course_code]["levels"].append("PGRAD")
                            if "Research" in coursetimetable:
                                courseinfo[currarea][course_code]["levels"].append("RESEARCH")

                            # load handbook page
                            coursepageraw = (get_html(handbookroot.format("undergraduate" if undergradcourse else "postgraduate") + course_code))
                            coursepage = str(coursepageraw)
                            #prereqs
                            courseinfo[currarea][course_code]["prereqs"] = []
                            if "Prerequisite" in coursepage:
                                with open("o.txt", "w") as f:
                                    f.write(coursepage)
                                m = re.search(r"Prerequisite:(.+?)<\/div>", coursepage)
                                if m is not None:
                                    m2 = re.findall(r"\s\w{4}\d{4}", str(m.group(0)))
                                    for prereqMatches in m2:
                                        courseinfo[currarea][course_code]["prereqs"].append(str(prereqMatches).strip())
                            
                            #excluded courses
                            # handbooksoup = bs(coursepageraw, 'html.parser')
                            # ex2 = handbooksoup.find(id="exclusion-rules")
                            # print(ex2)
                            # excluded = handbooksoup.find_all('div', attrs={"id":"exclusion-rules"})
                            # print("Excluded:")
                            # for ex in excluded:
                            #     pass
                            #     #print(ex)
                            # print("End excluded")

                            # Units of credits
                            if "6 Units of Credit" in coursepage:
                                courseinfo[currarea][course_code]["UOC"] = 6
                            elif "12 Units of Credit" in coursepage:
                                courseinfo[currarea][course_code]["UOC"] = 12

                            #Gen Ed
                            if "This course is offered as General Education" in coursepage:
                                courseinfo[currarea][course_code]["GENED"] = True
                            else:
                                courseinfo[currarea][course_code]["GENED"] = False


                            # class enrolment
                            courseinfo[currarea][course_code]["capacity"] = []
                            for i in range(4):
                                classutilpage = str(get_html(classutil_url.format(course_code[:4], utilterms[i])))
                                m = re.search(r'(name=\"{0}(.+?)(<td>).+?<td>.+?(.+?)(</td><td>).+?(</td><td>).+?(</td><td>)(.+?)(</td>).+?(cu00))'.format(course_code), classutilpage)
                                if m is not None:
                                    #print(m.group(0))
                                    #print(m.group(8))
                                    enrcap = (str(m.group(8))).split("/")
                                    courseinfo[currarea][course_code]["capacity"].append(int(enrcap[1]))
                            
                            courseinfo[currarea][course_code]["id"] = course_code
                            print(courseinfo[currarea][course_code])
                            
                            courses += 1
                            if courses > 1 and courses % 10 == 0:
                                with open("../../data/courseinfo.json", "w")  as f:
                                    json.dump(courseinfo, f)

                    time.sleep(1)
    
    print("Finished")
    with open("../../data/courseinfo.json", "w")  as f:
        json.dump(courseinfo, f)


if __name__ == '__main__':
    main()